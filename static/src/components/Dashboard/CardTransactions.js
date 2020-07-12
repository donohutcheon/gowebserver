import React, {useEffect, useState} from 'react';
import { lighten, makeStyles } from '@material-ui/core/styles';
import TableContainer from '@material-ui/core/TableContainer';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import {ERROR, SUCCESS} from "../../hooks/useApiRequest/actionTypes";
import useApiRequest from "../../hooks/useApiRequest/useApiRequest";
import TablePagination from "@material-ui/core/TablePagination";
import Paper from "@material-ui/core/Paper";
import EnhancedTableToolbar from "./EnhancedTableToolbar";
import EnhancedTableHead from "./EnhancedTableHead";

const headerCells = [
  { id: 'id', numeric: true, disablePadding: true, canSort: true, label: 'ID' },
  { id: 'amount', numeric: true, disablePadding: false, canSort: true, label: 'Amount' },
  { id: 'currency', numeric: false, disablePadding: false, canSort: false, label: 'Currency' },
  { id: 'dateTime', numeric: false, disablePadding: false, canSort: true, label: 'Time' },
  { id: 'category', numeric: false, disablePadding: false, canSort: false, label: 'Category' },
  { id: 'merchant', numeric: false, disablePadding: false, canSort: false, label: 'Merchant' },
  { id: 'city', numeric: false, disablePadding: false, canSort: false, label: 'City' },
  { id: 'country', numeric: false, disablePadding: false, canSort: false, label: 'Country' },
  { id: 'reference', numeric: false, disablePadding: false, canSort: false, label: 'Reference' },
];

/*const filterProps = [
  {id: 'amount', , onChange = ''  }
]*/

const useStyles = makeStyles((theme) => ({
  root: {
    width: '100%',
  },
  paper: {
    width: '100%',
    marginBottom: theme.spacing(2),
  },
  table: {
    minWidth: 750,
  },
  visuallyHidden: {
    border: 0,
    clip: 'rect(0 0 0 0)',
    height: 1,
    margin: -1,
    overflow: 'hidden',
    padding: 0,
    position: 'absolute',
    top: 20,
    width: 1,
  },
}));

export default function CardTransactions() {
  const classes = useStyles();
  const [cardTransactions, setCardTransactions] = useState([])
  const [order, setOrder] = React.useState('asc');
  const [orderBy, setOrderBy] = React.useState('id');
  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(10);
  const [nextIconButtonProps, setNextIconButtonProps] = React.useState({})

  const handleRequestSort = (event, property) => {
    const isAsc = orderBy === property && order === 'asc';
    setOrder(isAsc ? 'desc' : 'asc');
    setOrderBy(property);
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
    console.log("newPage: " + newPage)
    console.log("event: " + event)
    console.log("page: " + page)
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleRefreshTable = (event) => {
    makeRequest();
  };

  const [{ status, response }, makeRequest] = useApiRequest(
      `/api/me/card-transactions`,
      {
        verb: "get",
        params: {
          "page": page,
          "count": rowsPerPage,
          "sortField": orderBy,
          "sortDir": order,
        }
      }
  );
  useEffect(() => {
    console.log(process.env.NODE_ENV);
    makeRequest();
  }, [page, rowsPerPage, orderBy, order]);
  useEffect( () => {
    setNextIconButtonProps({disabled:cardTransactions.length < rowsPerPage})
  }, [cardTransactions])
  useEffect(() => {
    if(status === SUCCESS) {
      console.log("fetch card transactions", status, response.data)
      setCardTransactions(response.data.cardTransactions)
    } else if(status === ERROR) {
      console.log("fetch card transactions error", status, response.data)
    }
  }, [status]);
  return (
    <React.Fragment>
      <Paper className={classes.table} >
        <TableContainer>
          <EnhancedTableToolbar
              onRefreshTable={handleRefreshTable}
          />
          <Table size="small">
            <EnhancedTableHead
                classes={classes}
                order={order}
                orderBy={orderBy}
                onRequestSort={handleRequestSort}
                headCells={headerCells}
            />
            <TableBody>
              {cardTransactions.map((row) => (
                <TableRow key={row.id}>
                  <TableCell align="right">{row.id}</TableCell>
                  <TableCell align="right">{row.amount.value}</TableCell>
                  <TableCell>{row.currencyCode}</TableCell>
                  <TableCell>{row.dateTime}</TableCell>
                  <TableCell>{row.merchantCategoryName}</TableCell>
                  <TableCell>{row.merchantName}</TableCell>
                  <TableCell>{row.merchantCity}</TableCell>
                  <TableCell>{row.merchantCountryName}</TableCell>
                  <TableCell>{row.reference}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
        <TablePagination
            rowsPerPageOptions={[5, 10, 25]}
            component="div"
            count={-1}
            rowsPerPage={rowsPerPage}
            page={page}
            onChangePage={handleChangePage}
            onChangeRowsPerPage={handleChangeRowsPerPage}
            nextIconButtonProps={nextIconButtonProps}
        />
      </Paper>
    </React.Fragment>
  );
}
