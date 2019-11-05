import React, { Component } from "react";
import { Typography, Grid, Paper, FormControl, TextField, Button } from "@material-ui/core";

import ContentWrapper from "../../components/content-wrapper";

class Users extends Component {
  render() {
    return (
      <ContentWrapper align="top">
        <Paper style={{ padding: 24 }}>
          <Grid container xs={12} >

            {/* Formulario */}
            <Grid item xs={12} sm={9}>
              <Typography align="center" style={styles.textList}>Cadastrar novos usuários</Typography>
              <form>
                {/* Nome e Sobrenome */}
                <Grid container spacing={10}>
                  <Grid item xs={12} sm={6}>
                    {/* Nome */}
                    <FormControl fullWidth>
                      <TextField
                        id="input-nome"
                        label="Nome"
                        margin="normal"
                        variant="outlined"
                      />
                    </FormControl>
                  </Grid>
                  {/* Sobrenome */}
                  <Grid item xs={12} sm={6}>
                    <FormControl fullWidth>
                      <TextField
                        id="input-sobrenome"
                        label="RA"
                        margin="normal"
                        variant="outlined"
                      />
                    </FormControl>
                  </Grid>
                </Grid>
                {/* E-mail e Telefone */}
                <Grid container>
                  {/* E-mail */}
                  <Grid item xs={12} sm={6}>
                    <FormControl fullWidth>
                      <TextField
                        id="input-telefone"
                        label="Telefone"
                        margin="normal"
                        variant="outlined"
                      />
                    </FormControl>
                  </Grid>
                  <Grid item xs={12} sm={6}>
                    <FormControl fullWidth>
                      <TextField
                        id="input-placa"
                        label="Placa do veiculo"
                        margin="normal"
                        variant="outlined"
                      />
                    </FormControl>
                  </Grid>
                </Grid>
              </form>
              <Button variant="contained" color="primary">
                Adicionar usuário
      </Button>
            </Grid>

          </Grid>
        </Paper>
      </ContentWrapper>
    );
  }
}

const styles = {
  textList: {
    fontWeight: "bold",
    justifyContent: "center"
  }
};

export default Users;
