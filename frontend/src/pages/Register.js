import React, { Component } from "react";
import { Redirect } from "react-router";
import { connect } from "react-redux";
import * as authActions from "../modules/auth";

/**
  註冊頁面。
*/
class Register extends Component {
  state = {
    username: "",
    email: "",
    password: ""
  };

  usernameChangeHandler(event) {
    this.setState({
      username: event.target.value
    });
  }

  emailChangeHandler(event) {
    this.setState({
      email: event.target.value
    });
  }

  passwordChangeHandler(event) {
    this.setState({
      password: event.target.value
    });
  }

  registerHandler(event) {
    event.preventDefault();
    this.props.onRegister(
      this.state.username,
      this.state.email,
      this.state.password
    );
  }

  render() {
    if (this.props.user) return <Redirect to="/" />;

    return (
      <div className="columns">
        <div className="column is-4" />
        <div className="column is-4">
          <div className="field">
            <p className="control has-icons-left">
              <input
                className="input"
                type="text"
                placeholder="帳號"
                value={this.state.username}
                onChange={event => this.usernameChangeHandler(event)}
              />
              <span className="icon is-small is-left">
                <i className="fas fas fa-user" />
              </span>
            </p>
          </div>
          <div className="field">
            <p className="control has-icons-left">
              <input
                className="input"
                type="email"
                placeholder="Email"
                value={this.state.email}
                onChange={event => this.emailChangeHandler(event)}
              />
              <span className="icon is-small is-left">
                <i className="fas fa-envelope" />
              </span>
            </p>
          </div>
          <div className="field">
            <p className="control has-icons-left">
              <input
                className="input"
                type="password"
                placeholder="Password"
                value={this.state.password}
                onChange={event => this.passwordChangeHandler(event)}
              />
              <span className="icon is-small is-left">
                <i className="fas fa-lock" />
              </span>
            </p>
          </div>
          <div className="field">
            <p className="control">
              <button
                className="button is-success"
                onClick={event => this.registerHandler(event)}
              >
                註冊
              </button>
            </p>
          </div>
        </div>
        <div className="column is-4" />
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
    onRegister: (username, email, password) =>
      dispatch(authActions.register(username, email, password))
  };
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Register);
