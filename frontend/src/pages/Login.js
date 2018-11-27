import React, { Component } from "react";
import { Redirect } from "react-router";
import { connect } from "react-redux";
import * as authActions from "../modules/auth";

/**
  登入頁面。
*/
class Login extends Component {
  state = {
    email: "",
    password: ""
  };

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

  passwordKeyPressHandler(event) {
    if (event.key !== "Enter") return;
    this.loginHandler(event);
  }

  loginHandler(event) {
    event.preventDefault();
    this.props.onLogin(this.state.email, this.state.password);
  }

  render() {
    if (this.props.user) return <Redirect to="/" />;

    return (
      <div className="columns">
        <div className="column is-4" />
        <div className="column is-4">
          <div className="field">
            <label className="label">Email</label>
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
            <label className="label">密碼</label>
            <p className="control has-icons-left">
              <input
                className="input"
                type="password"
                placeholder="密碼"
                value={this.state.password}
                onChange={event => this.passwordChangeHandler(event)}
                onKeyPress={event => this.passwordKeyPressHandler(event)}
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
                onClick={event => this.loginHandler(event)}
              >
                登入
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
    onLogin: (email, password) => dispatch(authActions.login(email, password))
  };
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Login);
