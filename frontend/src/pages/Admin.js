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
    paginationKey: Math.random(), // 用來觸發分頁重新 render 並查詢資料。
    disableUserModalActivate: false, // 是否顯示停用使用者對話框。
    disableUserID: null // 停用使用者的 id。
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
      .get("/api/admin/users", {
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

  // 開啟停用使用者對話框。
  openDisableUserModal(id) {
    this.setState(
      produce(draft => {
        draft.disableUserID = id;
        draft.disableUserModalActivate = true;
      })
    );
  }

  // 關閉停用使用者對話框。
  closeDisableUserModal() {
    this.setState(
      produce(draft => {
        draft.disableUserID = null;
        draft.disableUserModalActivate = false;
      })
    );
  }

  disableUserHandler() {
    const id = this.state.disableUserID;
    axios.post(`/api/admin/users/disable/${id}`).then(response => {
      this.setState(
        produce(draft => {
          const index = draft.users.findIndex(user => user.id === id);
          draft.users[index] = Object.assign(
            draft.users[index],
            response.data.user
          );
        })
      );
      this.closeDisableUserModal();
    });
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
            {user.isDisabled === 0 && user.id !== this.props.user.id ? (
              <button
                className="button is-primary"
                onClick={event => this.openDisableUserModal(user.id)}
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
        <nav className="level">
          <div className="level-left">
            <div className="level-item">
              <h1 className="title">使用者管理</h1>
            </div>
          </div>
          <div className="level-right">
            <div className="level-item">
              <div className="field has-addons">
                <div className="control">
                  <input
                    className="input"
                    type="text"
                    value={this.state.inputSearchUser}
                    onChange={event =>
                      this.inputSearchUserChangeHandler(event.target.value)
                    }
                    onKeyPress={event =>
                      this.inputSearchUserKeyPressHandler(event)
                    }
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
            </div>
          </div>
        </nav>
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
        <div
          className={
            this.state.disableUserModalActivate ? "modal is-active" : "modal"
          }
        >
          <div className="modal-background" />
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="modal-card-title">停用使用者</p>
              <button
                className="delete"
                aria-label="close"
                onClick={event => this.closeDisableUserModal()}
              />
            </header>
            <section className="modal-card-body">確定停用使用者？</section>
            <footer className="modal-card-foot">
              <button
                className="button is-primary"
                onClick={event => this.disableUserHandler()}
              >
                確定
              </button>
              <button
                className="button"
                onClick={event => this.closeDisableUserModal()}
              >
                取消
              </button>
            </footer>
          </div>
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

export default connect(
  mapStateToProps,
  null
)(Admin);
