import React, { Component } from "react";
import { connect } from "react-redux";
import classes from "./spinner.module.css";

class Spinner extends Component {
  render() {
    return (
      <div className={this.props.loading ? "modal is-active" : "modal"}>
        <div className="modal-background" />
        <div
          className="modal-content has-text-centered"
          style={{ height: "100px", width: "100px" }}
        >
          <div className={classes.LdsRoller}>
            <div />
            <div />
            <div />
            <div />
            <div />
            <div />
            <div />
            <div />
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    loading: state.load.loading
  };
};

export default connect(
  mapStateToProps,
  null
)(Spinner);
