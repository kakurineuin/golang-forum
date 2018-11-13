import React, { Component } from 'react';
import { Link, Route, Switch, Redirect, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import Register from './pages/Register';
import Login from './pages/Login';
import Home from './pages/Home';
import * as authActions from './modules/auth';
import * as messageActions from './modules/message';
import './App.css';
import logo from './assets/logo.png';

class App extends Component {
  componentDidMount() {
    // 檢查是否登入過，若有就自動登入。
    this.props.onTryAutoLogin();
  }

  deleteMessageHandler(id) {
    this.props.onDeleteMessage(id);
  }

  logoutHandler() {
    this.props.onLogout();
  }

  render() {
    // 檢查有無訊息，若有就顯示。
    let notifications = [];
    const messageList = this.props.messageList;

    if (messageList && messageList.length > 0) {
      for (const message of messageList) {
        const { id, isError, text } = message;
        const classes = ['notification'];

        if (isError) {
          classes.push('is-danger');
        } else {
          classes.push('is-primary');
        }

        notifications.push(
          <div className={classes.join(' ')}>
            <button
              className="delete"
              onClick={event => this.deleteMessageHandler(id)}></button>
            {text}
          </div>
        );
      }
    }

    let myConfigs = null;
    let buttons = null;

    if (this.props.user) {
      myConfigs = <Link to="/my/configs" className="navbar-item">{this.props.user.account}</Link>;
      buttons = (
        <div className="buttons">
          <a className="button is-light"
            onClick={event => this.logoutHandler(event)}>登出</a>
        </div>
      );
    } else {
      buttons = (
        <div className="buttons">
          <Link to="/register" className="button is-primary">註冊</Link>
          <Link to="/login" className="button is-light">登入</Link>
        </div>
      );
    }

    return (
      <div>
        <nav className="navbar" role="navigation" aria-label="main navigation">
          <div className="navbar-brand">
            <Link to="/" className="navbar-item">
              <img src={logo} style={{ width: '100px', maxHeight: '52px' }} />
            </Link>
          </div>

          <div id="navbarBasicExample" className="navbar-menu">
            <div className="navbar-start">
              <Link to="/" className="navbar-item">Home</Link>
            </div>
            <div className="navbar-end">
              {myConfigs}
              <div className="navbar-item">
                {buttons}
              </div>
            </div>
          </div>
        </nav>

        {notifications}

        <div className="columns">
          <div className="column is-12">
            <Switch>
              <Route path="/register" component={Register} />
              <Route path="/login" component={Login} />
              <Route path="/" exact component={Home} />
              <Redirect to="/" />
            </Switch>
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    messageList: state.message.list,
    user: state.auth.user
  };
};

const mapDispatchToProps = dispatch => {
  return {
    onTryAutoLogin: () => dispatch(authActions.authCheckState()),
    onDeleteMessage: (id) => dispatch(messageActions.deleteMessage(id)),
    onLogout: () => dispatch(authActions.logout())
  };
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(App));
