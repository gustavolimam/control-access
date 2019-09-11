import React, { Component } from "react";
import { Typography, Grid, Paper } from "@material-ui/core";

import ContentWrapper from "../../components/content-wrapper";

class Users extends Component {
  render() {
    return (
      <ContentWrapper align="top">
        <Paper style={{ padding: 24 }}>
          <Grid container xs={12} spacing={24}>
            <Typography style={styles.textList}>Usuários</Typography>
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

export default Users;
