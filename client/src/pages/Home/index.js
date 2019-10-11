import React, { Component } from 'react';
import { Typography, Grid, Paper, Card } from '@material-ui/core';
import { PieChart, Pie, Cell, Legend } from 'recharts';
// import { info } from "../../services/firebase";

import ContentWrapper from '../../components/content-wrapper';
import TablesEvent from '../../components/tables';

// console.log(info);

// Constantes para testes do gráfico, futuramente será o retorno de uma API
const data0 = [
  { name: 'Livres', value: 800 },
  { name: 'Ocupadas', value: 200 }
];
const data1 = [
  { name: 'Livres', value: 300 },
  { name: 'Ocupadas', value: 700 }
];

// Constante que define os padrões de cores utilizados nos gráficos do dashboard
const COLORS = ['#00AEEF', '#EB060D'];

class Home extends Component {
  render() {
    return (
      <ContentWrapper align="top">
        <Paper style={{ padding: 24 }}>
          <Grid container xs={12} spacing={24}>
            <Grid container xs={12} justify="center">
              <Typography variant="h4" style={styles.textList}>
                Página Inicial
              </Typography>
            </Grid>
            <Grid item xs={6}>
              <Card>
                <Typography style={styles.textCard} variant="h6">
                  Campus
                </Typography>
                {/* Gráfico de vagas disponíveis no campus */}
                <PieChart
                  width={550}
                  height={250}
                  onMouseEnter={this.onPieEnter}
                >
                  <Pie
                    data={data0}
                    cx={275}
                    cy={100}
                    innerRadius={50}
                    outerRadius={100}
                    fill="#8884d8"
                    paddingAngle={5}
                    dataKey="value"
                  >
                    {data0.map((entry, index) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={COLORS[index % COLORS.length]}
                      />
                    ))}
                  </Pie>
                  <Legend
                    style={styles.textLegend}
                    layout="vertical"
                    align="right"
                  />
                </PieChart>
              </Card>
            </Grid>
            <Grid item item xs={6}>
              <Card>
                <Typography style={styles.textCard} variant="h6">
                  Bolsão
                </Typography>
                {/* Gráfico de vagas disponíveis no bolsão dos professores */}
                <PieChart
                  width={550}
                  height={250}
                  onMouseEnter={this.onPieEnter}
                >
                  <Pie
                    data={data1}
                    cx={275}
                    cy={100}
                    innerRadius={50}
                    outerRadius={100}
                    fill="#8884d8"
                    paddingAngle={5}
                    dataKey="value"
                  >
                    {data0.map((entry, index) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={COLORS[index % COLORS.length]}
                      />
                    ))}
                  </Pie>
                  <Legend
                    style={styles.textLegend}
                    layout="vertical"
                    align="right"
                  />
                </PieChart>
              </Card>
            </Grid>
          </Grid>
          <Grid container xs={12} spacing={24}>
            <TablesEvent />
          </Grid>
        </Paper>
      </ContentWrapper>
    );
  }
}

const styles = {
  textCard: {
    fontWeight: 'bold',
    textAlign: 'center'
  },
  textLegend: {
    fontWeight: 'bold'
  },
  textList: {
    fontWeight: 'bold'
  }
};

export default Home;
