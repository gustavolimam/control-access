import React, { Component } from "react";
import { BrowserRouter, Route } from "react-router-dom";

import Home from "./pages/Home/";
import Users from "./pages/Users/";
import Cameras from "./pages/Cameras/";
import AppHeader from "./components/AppBar";

export default class src extends Component {
  render() {
    return (
      <BrowserRouter>
        <AppHeader />
        <Route exact path="/" component={Home} />
        <Route path="/users" component={Users} />
        <Route path="/cams" component={Cameras} />
      </BrowserRouter>
    );
  }
}
