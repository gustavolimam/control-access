import React, { Component } from 'react';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Logo from '../assets/logoFacens.jpg';
import Button from '@material-ui/core/Button';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import { Grid, Typography } from '@material-ui/core';
import { Link } from 'react-router-dom';

class AppHeader extends Component {
  state = {
    top: false,
    left: false,
    bottom: false,
    right: false
  };

  toggleDrawer = (side, open) => () => {
    this.setState({
      [side]: open
    });
  };

  render() {
    const menu = [
      {
        name: 'Página Inicial',
        url: '/'
      },
      {
        name: 'Câmeras',
        url: '/cameras'
      },
      {
        name: 'Base de Usuários',
        url: '/users'
      }
    ];
    // Define a lista de diretórios possiveis na página web do Dashboard
    const menuList = (
      <div className={styles.list}>
        <List>
          {menu.map(text => (
            <ListItem button component={Link} to={text.url} key={text.name}>
              <ListItemText primary={text.name} />
            </ListItem>
          ))}
        </List>
      </div>
    );

    return (
      <AppBar position="fixed">
        <Toolbar style={styles.toolBar}>
          <img alt="Logo" src={Logo} width="80" />
          <Grid container justify="flex-end">
            <Button onClick={this.toggleDrawer('right', true)}>
              <Typography variant="h6" style={styles.textBar}>
                Controle de Acesso
              </Typography>
            </Button>
            <Drawer
              anchor="right"
              open={this.state.right}
              onClose={this.toggleDrawer('right', false)}
            >
              <div
                tabIndex={0}
                role="button"
                onClick={this.toggleDrawer('right', false)}
                onKeyDown={this.toggleDrawer('right', false)}
              >
                {menuList}
              </div>
            </Drawer>
          </Grid>
        </Toolbar>
      </AppBar>
    );
  }
}

const styles = {
  toolBar: {
    backgroundColor: '#00AEEF'
  },
  textBar: {
    color: 'white',
    fontWeight: 'bold'
  },
  list: {
    width: 250
  }
};

export default AppHeader;
