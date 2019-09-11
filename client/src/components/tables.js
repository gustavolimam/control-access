import React from "react";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Paper from "@material-ui/core/Paper";

function createData(name, placa, horaEntrada, horaSaida) {
  return { name, placa, horaEntrada, horaSaida };
}

const rows = [
  createData("Usuário 1", "EXI-7254", "06/05/19 15:59:22", "06/05/19 18:49:22"),
  createData("Usuário 2", "ABC-1234", "06/05/19 15:44:32", "06/05/19 19:59:22"),
  createData("Usuário 3", "ABC-4321", "06/05/19 14:50:13", "06/05/19 15:59:22"),
  createData("Usuário 4", "CBA-1234", "06/05/19 14:36:53", "06/05/19 17:59:22"),
  createData("Usuário 5", "CBA-4321", "06/05/19 14:33:30", "06/05/19 22:59:22")
];

function SimpleTable() {
  return (
    <Paper style={styles.root}>
      <Table style={styles.table}>
        <TableHead>
          <TableRow>
            <TableCell style={styles.tbl}>Usuários</TableCell>
            <TableCell align="right">Placa</TableCell>
            <TableCell align="right">Horário de Entrada</TableCell>
            <TableCell align="right">Horário de Saída</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map(row => (
            <TableRow key={row.name}>
              <TableCell component="th" scope="row">
                {row.name}
              </TableCell>
              <TableCell align="right">{row.placa}</TableCell>
              <TableCell align="right">{row.horaEntrada}</TableCell>
              <TableCell align="right">{row.horaSaida}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Paper>
  );
}

const styles = {
  root: {
    width: "100%",
    marginTop: 30,
    overflowX: "auto"
  },
  table: {
    minWidth: 650
  },
  tbl: {
    fontWeight: "bold"
  }
};
export default SimpleTable;
