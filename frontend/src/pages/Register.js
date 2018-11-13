import React, { Component } from 'react';
import { Redirect } from 'react-router'
import { connect } from 'react-redux';
import * as authActions from '../modules/auth'

class Register extends Component {
  state = {
    account: '',
    email: '',
    password: ''
  }

  accountChangedHandler(event) {
    this.setState({
      account: event.target.value
    });
  }

  emailChangedHandler(event) {
    this.setState({
      email: event.target.value
    });
  }

  passwordChangedHandler(event) {
    this.setState({
      password: event.target.value
    });
  }

  registerHandler(event) {
    event.preventDefault();
    this.props.onRegister(this.state.account, this.state.email, this.state.password);
  }

  render() {
    if (this.props.user) return <Redirect to="/" />;

    return (
      <div className="columns">
        <div className="column is-4">
        </div>
        <div className="column is-4">
          <div className="field">
            <p className="control has-icons-left">
              <input className="input"
                type="text"
                placeholder="帳號"
                value={this.state.account}
                onChange={event => this.accountChangedHandler(event)} />
              <span className="icon is-small is-left">
                <i className="fas fas fa-user"></i>
              </span>
            </p>
          </div>
          <div className="field">
            <p className="control has-icons-left">
              <input className="input"
                type="email"
                placeholder="Email"
                value={this.state.email}
                onChange={event => this.emailChangedHandler(event)} />
              <span className="icon is-small is-left">
                <i className="fas fa-envelope"></i>
              </span>
            </p>
          </div>
          <div className="field">
            <p className="control has-icons-left">
              <input className="input"
                type="password"
                placeholder="Password"
                value={this.state.password}
                onChange={event => this.passwordChangedHandler(event)} />
              <span className="icon is-small is-left">
                <i className="fas fa-lock"></i>
              </span>
            </p>
          </div>
          <div className="field">
            <p className="control">
              <button className="button is-success"
                onClick={event => this.registerHandler(event)}>
                註冊
              </button>
            </p>
          </div>
        </div>
        <div className="column is-4">
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    user: state.auth.user
  };
};

const mapDispatchToProps = dispatch => {
  return {
    onRegister: (account, email, password) =>
      dispatch(authActions.register(account, email, password))
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(Register);