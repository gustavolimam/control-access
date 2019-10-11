import React, { Component } from 'react';
import { Grid, Paper, Typography, CardMedia } from '@material-ui/core';
import Rodovia1 from '../../assets/rodovia1.jpg';
import Rodovia2 from '../../assets/rodovia2.jpg';
import Rodovia3 from '../../assets/rodovia3.jpg';
import Rodovia4 from '../../assets/rodovia4.jpg';

import ContentWrapper from '../../components/content-wrapper';

class Cameras extends Component {
  render() {
    return (
      <ContentWrapper align="top">
        <Paper style={{ padding: 24 }}>
          <Grid container xs={12} spacing={24} justify="center">
            <Typography variant="h4" style={styles.textList}>
              Câmeras
            </Typography>
            <Grid container xs={12} justify="center">
              <Grid xs={5}>
                <img src={Rodovia1} height={300} />
                <Typography style={{ textAlign: 'center' }}>
                  Câmera 1
                </Typography>
              </Grid>
              <Grid xs={5}>
                <img src={Rodovia2} height={300} />
                <Typography style={{ textAlign: 'center' }}>
                  Câmera 2
                </Typography>
              </Grid>
              <Grid xs={5}>
                <img src={Rodovia3} height={300} />
                <Typography style={{ textAlign: 'center' }}>
                  Câmera 3
                </Typography>
              </Grid>
              <Grid xs={5}>
                <img src={Rodovia4} height={300} />
                <Typography style={{ textAlign: 'center' }}>
                  Câmera 4
                </Typography>
              </Grid>
            </Grid>
          </Grid>
        </Paper>
      </ContentWrapper>
    );
  }
}

const styles = {
  textList: {
    fontWeight: 'bold'
  },
  media: {
    height: 140
  }
};

export default Cameras;
