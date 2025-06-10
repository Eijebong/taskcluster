import React, { Component } from 'react';
import { string, shape, func, arrayOf } from 'prop-types';
import { pipe, map, sort as rSort } from 'ramda';
import { camelCase } from 'camel-case';
import { withStyles } from '@material-ui/core/styles';
import TableRow from '@material-ui/core/TableRow';
import TableCell from '@material-ui/core/TableCell';
import DeleteIcon from 'mdi-react/DeleteIcon';
import { memoize } from '../../utils/memoize';
import Button from '../Button';
import sort from '../../utils/sort';
import ConnectionDataTable from '../ConnectionDataTable';
import { VIEW_ROLES_PAGE_SIZE } from '../../utils/constants';
import { pageInfo, role } from '../../utils/prop-types';
import Link from '../../utils/Link';

const sorted = pipe(
  rSort((a, b) => sort(a.node.roleId, b.node.roleId)),
  map(({ node: { roleId } }) => roleId)
);
const tableHeaders = ['Role ID'];
const iconSize = 16;

@withStyles(theme => ({
  roleIdLink: {
    display: 'flex',
    flexGrow: 1,
  },
  roleContainer: {
    paddingTop: theme.spacing(1.5),
    paddingBottom: theme.spacing(1.5),
    paddingLeft: theme.spacing(2),
    paddingRight: theme.spacing(2),
    display: 'flex',
    width: '100%',
  },
}))
export default class RolesTable extends Component {
  static propTypes = {
    rolesConnection: shape({
      edges: arrayOf(role),
      pageInfo,
    }).isRequired,
    onPageChange: func.isRequired,
    /** A search term to refine the list of roles */
    searchTerm: string,
    onDialogActionOpen: func.isRequired,
  };

  static defaultProps = {
    searchTerm: null,
  };

  state = {
    sortBy: tableHeaders[0],
    sortDirection: 'asc',
  };

  createSortedRolesConnection = memoize(
    (rolesConnection, sortBy, sortDirection) => {
      const sortByProperty = sortBy ? camelCase(sortBy) : '';

      if (!sortBy) {
        return rolesConnection;
      }

      return {
        ...rolesConnection,
        edges: [...rolesConnection.edges].sort((a, b) => {
          const firstElement =
            sortDirection === 'desc'
              ? b.node[sortByProperty]
              : a.node[sortByProperty];
          const secondElement =
            sortDirection === 'desc'
              ? a.node[sortByProperty]
              : b.node[sortByProperty];

          return sort(firstElement, secondElement);
        }),
      };
    },
    {
      serializer: ([rolesConnection, sortBy, sortDirection]) => {
        const ids = sorted(rolesConnection.edges);

        return `${ids.join('-')}-${sortBy}-${sortDirection}`;
      },
    }
  );

  handleHeaderClick = sortBy => {
    const toggled = this.state.sortDirection === 'desc' ? 'asc' : 'desc';
    const sortDirection = this.state.sortBy === sortBy ? toggled : 'desc';

    this.setState({ sortBy, sortDirection });
  };

  render() {
    const {
      classes,
      onPageChange,
      rolesConnection,
      searchTerm,
      onDialogActionOpen,
    } = this.props;
    const { sortBy, sortDirection } = this.state;
    const sortedRolesConnection = this.createSortedRolesConnection(
      rolesConnection,
      sortBy,
      sortDirection
    );

    return (
      <ConnectionDataTable
        searchTerm={searchTerm}
        connection={sortedRolesConnection}
        pageSize={VIEW_ROLES_PAGE_SIZE}
        onHeaderClick={this.handleHeaderClick}
        onPageChange={onPageChange}
        headers={tableHeaders}
        sortByHeader={sortBy}
        sortDirection={sortDirection}
        allowFilter
        filterFunc={({ node: role }, filterValue) =>
          String(role.roleId).includes(filterValue)
        }
        lazyRender
        renderRow={({ node: role }, style, key) => (
          <TableRow key={key || role.roleId} style={style} hover>
            <TableCell className={classes.roleContainer}>
              <Link
                className={classes.roleIdLink}
                to={`/auth/roles/${encodeURIComponent(role.roleId)}`}>
                {role.roleId}
              </Link>
              <Button
                requiresAuth
                tooltipProps={{ title: 'Delete Role' }}
                size="small"
                onClick={() => onDialogActionOpen(role.roleId)}>
                <DeleteIcon size={iconSize} />
              </Button>
            </TableCell>
          </TableRow>
        )}
      />
    );
  }
}
