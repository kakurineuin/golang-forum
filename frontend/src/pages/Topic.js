import React, { Component } from "react";
import PostEditor from "../components/PostEditor";
import Pagination from "../components/Pagination";
import { connect } from "react-redux";
import axios from "axios";
import produce from "immer";
import dateFns from "date-fns";

/**
  討論串頁面。
*/
class Topic extends Component {
  state = {
    posts: [], // 全部文章。
    totalCount: 0, // 文章總數。
    content: "", // 新增回覆內文。
    paginationKey: Math.random(), // 用來觸發分頁重新 render 並查詢資料。
    postOnUpdate: {
      content: ""
    }, // 修改的文章。
    updatePostModalActivate: false, // 是否顯示修改文章對話框。
    postOnDelete: {}, // 刪除的文章。
    deletePostModalActivate: false // 是否顯示刪除文章對話框。
  };

  // 新增回覆內文。
  contentChangeHandler(value) {
    this.setState(
      produce(draft => {
        draft.content = value;
      })
    );
  }

  // 修改文章。
  postOnUpdateChangeHandler(value) {
    this.setState(
      produce(draft => {
        draft.postOnUpdate.content = value;
      })
    );
  }

  // 新增回覆。
  createReplyHandler() {
    console.log("props", this.props);
    console.log("state", this.state);
    axios
      .post(`/api/topics/${this.props.match.params.category}`, {
        userProfileId: this.props.user.id,
        topic: this.state.posts[0].topic,
        content: this.state.content,
        replyPostId: parseInt(this.props.match.params.id, 10)
      })
      .then(response => {
        console.log("create reply response", response);
        this.setState(
          produce(draft => {
            draft.content = "";
            draft.paginationKey = Math.random();
          })
        );
      });
  }

  // 開啟修改文章對話框。
  openUpdatePostModal(post) {
    this.setState(
      produce(draft => {
        draft.postOnUpdate = Object.assign({}, post);
        draft.updatePostModalActivate = true;
      })
    );
  }

  // 關閉修改文章對話框。
  closeUpdatePostModal() {
    this.setState(
      produce(draft => {
        draft.postOnUpdate = { content: "" };
        draft.updatePostModalActivate = false;
      })
    );
  }

  // 修改文章。
  updatePostHandler() {
    console.log("props", this.props);
    console.log("state", this.state);
    const category = this.props.match.params.category;
    const postId = this.state.postOnUpdate.id;
    axios
      .put(`/api/topics/${category}/${postId}`, {
        content: this.state.postOnUpdate.content
      })
      .then(response => {
        console.log("update post response", response);
        const updatedPost = response.data.post;
        this.setState(
          produce(draft => {
            const index = draft.posts.findIndex(
              post => post.id === updatedPost.id
            );
            draft.posts[index] = Object.assign(draft.posts[index], updatedPost);
            draft.postOnUpdate = { content: "" };
            draft.updatePostModalActivate = false;
          })
        );
      });
  }

  // 開啟刪除文章對話框。
  openDeletePostModal(post) {
    this.setState(
      produce(draft => {
        draft.postOnDelete = Object.assign({}, post);
        draft.deletePostModalActivate = true;
      })
    );
  }

  // 關閉刪除文章對話框。
  closeDeletePostModal() {
    this.setState(
      produce(draft => {
        draft.postOnDelete = {};
        draft.deletePostModalActivate = false;
      })
    );
  }

  // 刪除文章。
  deletePostHandler() {
    console.log("props", this.props);
    console.log("state", this.state);
    const category = this.props.match.params.category;
    const postId = this.state.postOnDelete.id;
    axios.delete(`/api/topics/${category}/${postId}`).then(response => {
      console.log("delete post response", response);
      const deletedPost = response.data.post;
      this.setState(
        produce(draft => {
          const index = draft.posts.findIndex(
            post => post.id === deletedPost.id
          );
          draft.posts[index] = Object.assign(draft.posts[index], deletedPost);
          draft.postOnDelete = {};
          draft.deletePostModalActivate = false;
        })
      );
    });
  }

  // 查詢此主題文章。
  findPostsByTopicId(offset, limit) {
    const category = this.props.match.params.category;
    const id = this.props.match.params.id;
    axios
      .get(`/api/topics/${category}/${id}`, {
        params: {
          offset,
          limit
        }
      })
      .then(response => {
        console.log(response);
        this.setState(
          produce(draft => {
            draft.posts = response.data.posts;
            draft.totalCount = response.data.totalCount;
          })
        );
      });
  }

  goBackHandler() {
    this.props.history.goBack();
  }

  render() {
    let createReply = null;

    if (this.props.user) {
      createReply = (
        <div className="box">
          <div className="field">
            <label className="label">回覆</label>
            <div className="control">
              <PostEditor
                value={this.state.content}
                changed={value => this.contentChangeHandler(value)}
              />
            </div>
          </div>
          <div className="field">
            <div className="control">
              <button
                className="button is-primary"
                onClick={event => this.createReplyHandler()}
              >
                回覆
              </button>
            </div>
          </div>
        </div>
      );
    }

    const postId = parseInt(this.props.match.params.id, 10);
    const posts = this.state.posts.map((post, index) => {
      let updateButton = null;
      let deleteButton = null;
      const user = this.props.user;

      // 只能修改自己的文章。
      if (!post.deletedAt && user && post.username === user.username) {
        updateButton = (
          <button
            className="button is-primary"
            onClick={event => this.openUpdatePostModal(post)}
          >
            修改
          </button>
        );
      }

      // 只能刪除自己的文章，而系統管理員可以刪除每個人的文章。
      if (
        !post.deletedAt &&
        user &&
        (post.username === user.username || user.role === "admin")
      ) {
        deleteButton = (
          <button
            className="button"
            onClick={event => this.openDeletePostModal(post)}
          >
            刪除
          </button>
        );
      }

      let content = null;
      let postContent = post.content;

      // 若文章已刪除。
      if (post.deletedAt) {
        // 刪除時間。
        const deletedAt = dateFns.format(
          new Date(post.deletedAt),
          "YYYY/MM/DD HH:mm:ss"
        );
        postContent = deletedAt + " - " + postContent;
      }

      // 若是主題文章，顯示主題。
      content = (
        <td>
          {post.id === postId ? (
            <div>
              <h1 className="title">{post.topic}</h1>
              <hr />
            </div>
          ) : null}
          <div className="ql-snow">
            <div
              className="ql-editor"
              dangerouslySetInnerHTML={{ __html: postContent }}
            />
          </div>
          {updateButton}
          {deleteButton}
        </td>
      );

      // 新增時間。
      const createdAt = dateFns.format(
        new Date(post.createdAt),
        "YYYY/MM/DD HH:mm:ss"
      );

      // 修改時間。
      const updatedAt = dateFns.format(
        new Date(post.updatedAt),
        "YYYY/MM/DD HH:mm:ss"
      );

      return (
        <tr key={post.id}>
          <td>
            <h6 className="title is-6">{post.username}</h6>
            <div className="is-size-7">{createdAt + " 新增"}</div>
            {updatedAt !== createdAt ? (
              <div className="is-size-7">{updatedAt + " 修改"}</div>
            ) : null}
            <br />
          </td>
          {content}
        </tr>
      );
    });

    return (
      <div className="container">
        <table className="table is-bordered is-striped is-fullwidth">
          <thead>
            <tr>
              <th width="170px">作者</th>
              <th>文章</th>
            </tr>
          </thead>
          <tbody>{posts}</tbody>
        </table>
        <Pagination
          key={this.state.paginationKey}
          totalCount={this.state.totalCount}
          findData={(offset, limit) => this.findPostsByTopicId(offset, limit)}
        />
        <br />
        <nav className="level">
          <div className="level-left">
            <div className="level-item" />
          </div>
          <div className="level-right">
            <div className="level-item">
              <button
                className="button"
                onClick={event => this.goBackHandler()}
              >
                回上一頁
              </button>
            </div>
          </div>
        </nav>
        {createReply}
        <div
          className={
            this.state.updatePostModalActivate ? "modal is-active" : "modal"
          }
        >
          <div className="modal-background" />
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="modal-card-title">修改文章</p>
              <button
                className="delete"
                aria-label="close"
                onClick={event => this.closeUpdatePostModal()}
              />
            </header>
            <section className="modal-card-body">
              <PostEditor
                value={
                  this.state.postOnUpdate ? this.state.postOnUpdate.content : ""
                }
                changed={value => this.postOnUpdateChangeHandler(value)}
              />
            </section>
            <footer className="modal-card-foot">
              <button
                className="button is-primary"
                onClick={event => this.updatePostHandler()}
              >
                確定
              </button>
              <button
                className="button"
                onClick={event => this.closeUpdatePostModal()}
              >
                取消
              </button>
            </footer>
          </div>
        </div>
        <div
          className={
            this.state.deletePostModalActivate ? "modal is-active" : "modal"
          }
        >
          <div className="modal-background" />
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="modal-card-title">刪除文章</p>
              <button
                className="delete"
                aria-label="close"
                onClick={event => this.closeDeletePostModal()}
              />
            </header>
            <section className="modal-card-body">確定刪除文章？</section>
            <footer className="modal-card-foot">
              <button
                className="button is-primary"
                onClick={event => this.deletePostHandler()}
              >
                確定
              </button>
              <button
                className="button"
                onClick={event => this.closeDeletePostModal()}
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
)(Topic);
