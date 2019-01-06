import React, { Component } from "react";
import { connect } from "react-redux";
import classes from "./spinner.css";

class Spinner extends Component {
  render() {
    return (
      // <div className={this.props.loading ? "modal is-active" : "modal"}>
      //   <div className="modal-background" />
      //   <div className="modal-content">
      <div className={classes.Loader}>Loading...</div>
      //   </div>
      // </div>
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
