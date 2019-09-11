import React, { Component } from "react";
import { Grid, Paper, Typography } from "@material-ui/core";

import ContentWrapper from "../../components/content-wrapper";

class Cameras extends Component {
  render() {
    return (
      <ContentWrapper align="top">
        <Paper style={{ padding: 24 }}>
          <Grid container xs={12} spacing={24}>
            <Typography style={styles.textList}>Câmeras</Typography>
          </Grid>
        </Paper>
      </ContentWrapper>
    );
  }
}

const styles = {
  textList: {
    fontWeight: "bold"
  }
};

export default Cameras;
