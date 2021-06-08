import React from 'react'
import { makeStyles, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@material-ui/core'

const useStyles = makeStyles({
  table: {
    minWidth: 650,
  },
})

export default function LogTable(): JSX.Element {
  const classes = useStyles()

  const rows = [
    {key: 1, timestamp: '123456789', log: 'My log message'},
    {key: 2, timestamp: '123456789', log: 'My log message'},
    {key: 3, timestamp: '123456789', log: 'My log message'},
    {key: 4, timestamp: '123456789', log: 'My log message'},
    {key: 5, timestamp: '123456789', log: 'My log message'},
    {key: 6, timestamp: '123456789', log: 'My log message'},
    {key: 7, timestamp: '123456789', log: 'My log message'},
  ]
  return (
    <TableContainer>
      <Table className={classes.table} aria-label="simple table">
        <TableHead >
          <TableRow>
            <TableCell>Timestamp</TableCell>
            <TableCell>Log Message</TableCell>
          </TableRow>
        </TableHead>

        <TableBody>
          {rows.map((row) => (
            <TableRow key={row.key}> 
              <TableCell>{row.timestamp}</TableCell>
              <TableCell>{row.log}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}
