import React, { Component, Fragment } from 'react';
import { func, string } from 'prop-types';
import { parse, stringify } from 'qs';
import { withRouter } from 'react-router-dom';
import { formatDistanceStrict, parseISO } from 'date-fns';
import { pipe, map, sort as rSort } from 'ramda';
import memoize from 'fast-memoize';
import { withStyles } from '@material-ui/core/styles';
import { TableCell, TableRow, Typography } from '@material-ui/core';
import LinkIcon from 'mdi-react/LinkIcon';
import DeleteIcon from 'mdi-react/DeleteIcon';
import { WorkerManager } from 'taskcluster-client-web';
import Button from '../Button';
import CopyToClipboardTableCell from '../CopyToClipboardTableCell';
import StatusLabel from '../StatusLabel';
import DateDistance from '../DateDistance';
import TableCellItem from '../TableCellItem';
import ConnectionDataTable from '../ConnectionDataTable';
import Label from '../Label';
import DialogAction from '../DialogAction';
import { NULL_PROVIDER, VIEW_WORKERS_PAGE_SIZE } from '../../utils/constants';
import { workers } from '../../utils/prop-types';
import { withAuth } from '../../utils/Auth';
import Link from '../../utils/Link';
import sort from '../../utils/sort';

const sorted = pipe(
  rSort((a, b) => sort(a.node.workerId, b.node.workerId)),
  map(
    ({ node: { workerId, latestTask } }) =>
      `${workerId}.${latestTask?.run?.taskId ?? '-'}.${latestTask?.run?.runId ??
        '-'}`
  )
);

@withAuth
@withRouter
@withStyles(theme => ({
  linksIcon: {
    marginLeft: theme.spacing(1),
  },
  button: {
    marginLeft: -theme.spacing(2),
    marginRight: theme.spacing(1),
    borderRadius: 4,
  },
}))
/**
 * Display relevant information about workers in a table.
 */
export default class WorkersTable extends Component {
  static propTypes = {
    /** Workers GraphQL PageConnection instance. */
    workersConnection: workers,
    /** Callback function fired when a page is changed. */
    onPageChange: func.isRequired,
    /** Worker type name */
    workerType: string.isRequired,
    /** Provisioner identifier */
    provisionerId: string.isRequired,
  };

  static defaultProps = {
    /** Workers GraphQL PageConnection instance. */
    workersConnection: {
      edges: [],
      pageInfo: {},
    },
  };

  state = {
    error: null,
    open: false,
    title: '',
    body: '',
    confirmText: '',
    taskQueueId: '',
    workerGroup: '',
    workerId: '',
  };

  createSortedWorkersConnection = memoize(
    (workersConnection, sortBy, sortDirection) => {
      if (!sortBy) {
        return workersConnection;
      }

      return {
        ...workersConnection,
        edges: [...workersConnection.edges].sort((a, b) => {
          const firstElement =
            sortDirection === 'desc'
              ? this.valueFromNode(b.node)
              : this.valueFromNode(a.node);
          const secondElement =
            sortDirection === 'desc'
              ? this.valueFromNode(a.node)
              : this.valueFromNode(b.node);

          return sort(firstElement, secondElement);
        }),
      };
    },
    {
      serializer: ([workersConnections, sortBy, sortDirection]) => {
        const ids = sorted(workersConnections.edges);

        return `${ids.join('-')}-${sortBy}-${sortDirection}`;
      },
    }
  );

  handleDialogActionOpen = (taskQueueId, workerGroup, workerId) => () => {
    this.setState({
      open: true,
      title: 'Terminate Worker?',
      body: `This will terminate the worker with id ${workerId} in group ${workerGroup} within worker pool ${taskQueueId}.`,
      confirmText: 'Terminate Worker',
      taskQueueId,
      workerGroup,
      workerId,
    });
  };

  handleDeleteClick = async () => {
    const { taskQueueId, workerGroup, workerId } = this.state;
    const { user } = this.props;

    this.setState({
      error: null,
    });

    try {
      if (typeof user?.credentials === 'undefined') {
        throw new Error('User credentials not found. Please log in.');
      }

      const wm = new WorkerManager({
        rootUrl: window.env.TASKCLUSTER_ROOT_URL,
        credentials: user.credentials,
        authorizedScopes: [
          `worker-manager:remove-worker:${taskQueueId}/${workerGroup}/${workerId}`,
        ],
      });

      await wm.removeWorker(taskQueueId, workerGroup, workerId);
      this.setState({
        open: false,
      });
    } catch (error) {
      this.handleDialogActionError(error);
    }
  };

  handleDialogActionError = error => {
    this.setState({
      error,
    });
  };

  handleDialogActionClose = () => {
    this.setState({
      error: null,
      open: false,
    });
  };

  handleHeaderClick = sortByHeader => {
    const query = parse(this.props.location.search.slice(1));
    const sortBy = sortByHeader;
    const toggled = query.sortDirection === 'desc' ? 'asc' : 'desc';
    const sortDirection = query.sortBy === sortBy ? toggled : 'desc';

    query.sortBy = sortBy;
    query.sortDirection = sortDirection;
    this.props.history.replace({
      search: stringify(query, { addQueryPrefix: true }),
    });
  };

  valueFromNode(node) {
    const query = parse(this.props.location.search.slice(1));
    const mapping = {
      'Worker Group': node.workerGroup,
      'Worker ID': node.workerId,
      'Worker State': node.state,
      'Worker Capacity': node.capacity,
      'Last Active': node.lastDateActive,
      'First Claim': node.firstClaim,
      'Most Recent Task': node.latestTask?.run?.taskId,
      'Task State': node.latestTask?.run?.state,
      'Task Started': node.latestTask?.run?.started,
      'Task Resolved': node.latestTask?.run?.resolved,
      Quarantined: node.quarantineUntil,
    };

    return mapping[query.sortBy];
  }

  componentDidMount() {
    const query = parse(this.props.location.search.slice(1));

    if (query.sortBy) return;

    this.props.history.replace({
      search: stringify(
        { sortBy: 'Last Active', sortDirection: 'desc' },
        { addQueryPrefix: true }
      ),
    });
  }

  render() {
    const query = parse(this.props.location.search.slice(1));
    const { sortBy, sortDirection } = query.sortBy
      ? query
      : { sortBy: 'Last Active', sortDirection: 'desc' };
    const {
      provisionerId,
      workerType,
      onPageChange,
      workersConnection,
      classes,
      ...props
    } = this.props;
    const { open, error, title, confirmText, body } = this.state;
    const iconSize = 16;
    const connection = this.createSortedWorkersConnection(
      workersConnection,
      sortBy,
      sortDirection
    );

    return (
      <Fragment>
        <ConnectionDataTable
          connection={connection}
          pageSize={VIEW_WORKERS_PAGE_SIZE}
          sortByHeader={sortBy}
          sortDirection={sortDirection}
          onHeaderClick={this.handleHeaderClick}
          onPageChange={onPageChange}
          renderRow={({
            node: {
              workerId,
              workerGroup,
              latestTask,
              firstClaim,
              quarantineUntil,
              lastDateActive,
              state,
              capacity,
              providerId,
              taskQueueId,
            },
          }) => (
            <TableRow key={workerId}>
              <TableCell>{workerGroup}</TableCell>
              <TableCell>
                <Link
                  to={`/provisioners/${provisionerId}/worker-types/${workerType}/workers/${workerGroup}/${workerId}`}>
                  <TableCellItem button>
                    {workerId}
                    <LinkIcon className={classes.linksIcon} size={iconSize} />
                  </TableCellItem>
                </Link>
              </TableCell>
              <TableCell>
                {state ? (
                  <StatusLabel state={state.toUpperCase()} />
                ) : (
                  <em>n/a</em>
                )}
              </TableCell>
              <TableCell> {capacity || 0} </TableCell>
              {lastDateActive ? (
                <CopyToClipboardTableCell
                  tooltipTitle={lastDateActive}
                  textToCopy={lastDateActive}
                  text={<DateDistance from={lastDateActive} />}
                />
              ) : (
                <TableCell>
                  <em>n/a</em>
                </TableCell>
              )}
              <CopyToClipboardTableCell
                tooltipTitle={firstClaim}
                textToCopy={firstClaim}
                text={<DateDistance from={firstClaim} />}
              />
              <TableCell>
                {latestTask?.run ? (
                  <Link
                    to={`/tasks/${latestTask.run.taskId}/runs/${latestTask.run.runId}`}>
                    <TableCellItem button>
                      {latestTask.run.taskId}
                      <LinkIcon className={classes.linksIcon} size={iconSize} />
                    </TableCellItem>
                  </Link>
                ) : (
                  <em>n/a</em>
                )}
              </TableCell>
              <TableCell>
                {latestTask?.run ? (
                  <StatusLabel state={latestTask.run.state} />
                ) : (
                  <em>n/a</em>
                )}
              </TableCell>
              {latestTask?.run ? (
                <CopyToClipboardTableCell
                  tooltipTitle={latestTask.run.started}
                  textToCopy={latestTask.run.started}
                  text={<DateDistance from={latestTask.run.started} />}
                />
              ) : (
                <TableCell>n/a</TableCell>
              )}
              {latestTask?.run?.resolved ? (
                <CopyToClipboardTableCell
                  tooltipTitle={latestTask.run.resolved}
                  textToCopy={latestTask.run.resolved}
                  text={<DateDistance from={latestTask.run.resolved} />}
                />
              ) : (
                <TableCell>n/a</TableCell>
              )}
              <TableCell>
                {quarantineUntil ? (
                  formatDistanceStrict(new Date(), parseISO(quarantineUntil), {
                    unit: 'day',
                  })
                ) : (
                  <em>n/a</em>
                )}
              </TableCell>
              <TableCell>
                {providerId !== NULL_PROVIDER ? (
                  <Button
                    requiresAuth
                    variant="outlined"
                    endIcon={<DeleteIcon size={iconSize} />}
                    onClick={this.handleDialogActionOpen(
                      taskQueueId,
                      workerGroup,
                      workerId
                    )}
                    tooltipProps={{ title: '' }}>
                    Terminate
                  </Button>
                ) : (
                  <Label mini status="warning" className={classes.button}>
                    Scheduled for termination
                  </Label>
                )}
              </TableCell>
            </TableRow>
          )}
          headers={[
            'Worker Group',
            'Worker ID',
            'Worker State',
            'Worker Capacity',
            'Last Active',
            'First Claim',
            'Most Recent Task',
            'Task State',
            'Task Started',
            'Task Resolved',
            'Quarantined',
            '',
          ]}
          {...props}
        />
        <DialogAction
          open={open}
          onSubmit={this.handleDeleteClick}
          onClose={this.handleDialogActionClose}
          onError={this.handleDialogActionError}
          error={error}
          title={title}
          body={<Typography>{body}</Typography>}
          confirmText={confirmText}
        />
      </Fragment>
    );
  }
}
