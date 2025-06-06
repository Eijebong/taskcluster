package interactive

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	MsgStdin  = 1
	MsgResize = 2
)

type InteractiveJob struct {
	inner  InteractiveInnerType
	errors chan error
	done   chan struct{}
	wsLock sync.Mutex
	conn   *websocket.Conn
	ctx    context.Context
}

func CreateInteractiveJob(createCmd CreateInteractiveProcess, conn *websocket.Conn, ctx context.Context) (itj *InteractiveJob, err error) {
	itj = &InteractiveJob{
		// size of 3 is because there
		// are only ever 3 goroutines
		// who write to this channel
		// and we don't want to block
		errors: make(chan error, 3),
		done:   make(chan struct{}),
		wsLock: sync.Mutex{},
		conn:   conn,
		ctx:    ctx,
	}

	cmd, err := createCmd()
	if err != nil {
		itj.reportError(fmt.Sprintf("Error while getting command %v", err))
		return
	}
	itj.Setup(cmd)

	go itj.copyCommandOutputStream()
	go itj.handleWebsocketMessages()

	return itj, err
}

func (itj *InteractiveJob) Terminate() (err error) {
	select {
	case <-itj.done:
		return nil
	default:
		return itj.terminate()
	}
}

func (itj *InteractiveJob) copyCommandOutputStream() {
	buf := make([]byte, 4096)
	for {
		select {
		case <-itj.ctx.Done():
			return
		case <-itj.done:
			return
		default:
			n, err := itj.readPty(buf)
			if err != nil {
				if err == io.EOF {
					continue
				}
				return
			}
			if n == 0 {
				continue
			}
			if err := itj.writeWsMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				itj.errors <- err
				return
			}
		}
	}
}

func (itj *InteractiveJob) handleWebsocketMessages() {
	for {
		select {
		case <-itj.ctx.Done():
			return
		case <-itj.done:
			return
		case err := <-itj.errors:
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					itj.reportError(fmt.Sprintf("Error occured: %v", err))
				}
			}
			err = itj.Terminate()
			if err != nil {
				log.Printf("Error while terminating process: %v", err)
			}
			return
		default:
			_, msg, err := itj.conn.ReadMessage()
			if err != nil {
				itj.errors <- err
				continue
			}

			if len(msg) == 0 {
				return
			}

			switch msg[0] {
			case MsgStdin:
				if _, err := itj.writePty(msg[1:]); err != nil {
					itj.errors <- err
				}
			case MsgResize:
				width := binary.LittleEndian.Uint16(msg[1:3])
				height := binary.LittleEndian.Uint16(msg[3:])
				err = itj.resizePty(width, height)
				if err != nil {
					itj.errors <- err
				}
			default:
				log.Printf("Unknown message code received from interactive task")
			}
		}
	}
}

func (itj *InteractiveJob) reportError(errorMessage string) {
	log.Println(errorMessage)
	err := itj.writeWsMessage(websocket.BinaryMessage, []byte(errorMessage))
	if err != nil {
		log.Println("Error while reporting error to client")
	}
}

func (itj *InteractiveJob) writeWsMessage(messageType int, message []byte) (err error) {
	itj.wsLock.Lock()
	defer itj.wsLock.Unlock()
	return itj.conn.WriteMessage(messageType, message)
}
