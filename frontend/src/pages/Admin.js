import React, { Component } from "react";
import Pagination from "../components/Pagination";
import { connect } from "react-redux";
import axios from "axios";
import produce from "immer";
import dateFns from "date-fns";

/**
  系統管理員頁面，提供系統管理員停用使用者功能。
*/
class Admin extends Component {
  state = {
    users: [], // 使用者列表。
    totalCount: 0, // 主題總筆數。
    inputSearchUser: "", // 搜尋輸入框輸入的值。
    searchUser: "", // 請求參數：搜尋的使用者名稱。
    paginationKey: Math.random() // 用來觸發分頁重新 render 並查詢資料。
  };

  inputSearchUserChangeHandler(value) {
    this.setState(
      produce(draft => {
        draft.inputSearchUser = value;
      })
    );
  }

  inputSearchUserKeyPressHandler(event) {
    if (event.key !== "Enter") return;
    this.searchHandler();
  }

  searchHandler() {
    this.setState(
      produce(draft => {
        draft.searchUser = draft.inputSearchUser;
        draft.inputSearchUser = "";
        draft.paginationKey = Math.random(); // 觸發分頁重新查詢。
      })
    );
  }

  deleteTagHandler() {
    this.setState(
      produce(draft => {
        draft.searchUser = "";
        draft.paginationKey = Math.random(); // 觸發分頁重新查詢。
      })
    );
  }

  findUsers(offset, limit) {
    axios
      .get("/api/users", {
        params: {
          searchUser: this.state.searchUser,
          offset,
          limit
        }
      })
      .then(response => {
        this.setState(
          produce(draft => {
            draft.users = response.data.users;
            draft.totalCount = response.data.totalCount;
          })
        );
      });
  }

  disableUserHandler(id) {
    // TODO: 停用使用者。
  }

  render() {
    const users = this.state.users.map((user, index) => {
      return (
        <tr key={user.id}>
          <td>{user.username}</td>
          <td>{user.email}</td>
          <td>
            {dateFns.format(new Date(user.createdAt), "YYYY/MM/DD HH:mm:ss")}
          </td>
          <td>
            {user.id !== this.props.user.id ? (
              <button
                className="button is-primary"
                onClick={event => this.disableUserHandler(user.id)}
              >
                停用
              </button>
            ) : null}
          </td>
        </tr>
      );
    });

    // 顯示搜尋主題的標籤。
    let tagSearchUser = null;

    if (this.state.searchUser) {
      tagSearchUser = (
        <span className="tag is-medium">
          {this.state.searchUser}
          <button
            className="delete is-small"
            onClick={event => this.deleteTagHandler()}
          />
        </span>
      );
    }

    return (
      <div className="container">
        <h1 className="title">使用者管理</h1>
        <div className="field has-addons">
          <div className="control">
            <input
              className="input"
              type="text"
              value={this.state.inputSearchUser}
              onChange={event =>
                this.inputSearchUserChangeHandler(event.target.value)
              }
              onKeyPress={event => this.inputSearchUserKeyPressHandler(event)}
              placeholder="搜尋使用者名稱"
            />
          </div>
          <div className="control">
            <button
              className="button is-info"
              onClick={event => this.searchHandler()}
            >
              搜尋
            </button>
          </div>
        </div>
        {tagSearchUser}
        <hr />
        <table className="table is-bordered is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>使用者名稱</th>
              <th>Email</th>
              <th>建立時間</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>{users}</tbody>
        </table>
        <Pagination
          key={this.state.paginationKey}
          totalCount={this.state.totalCount}
          findData={(offset, limit) => this.findUsers(offset, limit)}
        />
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    user: state.auth.user
  };
};

export default connect(
  mapStateToProps,
  null
)(Admin);
