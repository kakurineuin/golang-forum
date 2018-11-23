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
    updatePostModalActivate: false // 是否顯示修改文章對話框。
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
        userProfileID: this.props.user.id,
        topic: this.state.posts[0].topic,
        content: this.state.content,
        replyPostID: parseInt(this.props.match.params.id, 10)
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

  // 修改文章。
  updatePostHandler() {
    console.log("props", this.props);
    console.log("state", this.state);
    const category = this.props.match.params.category;
    const postID = this.state.postOnUpdate.id;
    axios
      .put(`/api/topics/${category}/${postID}`, {
        content: this.state.postOnUpdate.content
      })
      .then(response => {
        console.log("update post response", response);
        const updatedPost = response.data.post;

        this.setState(
          produce(draft => {
            const post = draft.posts.find(post => post.id === updatedPost.id);
            post.content = updatedPost.content;
            post.updatedAt = updatedPost.updatedAt;
            draft.postOnUpdate = { content: "" };
            draft.updatePostModalActivate = false;
          })
        );
      });
  }

  // 查詢此主題文章。
  findPostsByTopicID(offset, limit) {
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

  openUpdatePostModal(post) {
    this.setState(
      produce(draft => {
        draft.postOnUpdate = Object.assign({}, post);
        draft.updatePostModalActivate = true;
      })
    );
  }

  closeUpdatePostModal() {
    this.setState(
      produce(draft => {
        draft.postOnUpdate = { content: "" };
        draft.updatePostModalActivate = false;
      })
    );
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

    const postID = parseInt(this.props.match.params.id, 10);
    const posts = this.state.posts.map((post, index) => {
      let updateButton = null;

      if (post.username === this.props.user.username) {
        updateButton = (
          <button
            className="button is-primary"
            onClick={event => this.openUpdatePostModal(post)}
          >
            修改
          </button>
        );
      }

      let content = null;

      // 若是主題文章，顯示主題。
      content = (
        <td>
          {post.id === postID ? (
            <div>
              <h1 className="title">{post.topic}</h1>
              <hr />
            </div>
          ) : null}
          <div className="ql-snow">
            <div
              className="ql-editor"
              dangerouslySetInnerHTML={{ __html: post.content }}
            />
          </div>
          {updateButton}
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
          findData={(offset, limit) => this.findPostsByTopicID(offset, limit)}
        />
        <br />
        {createReply}
        <div
          class={
            this.state.updatePostModalActivate ? "modal is-active" : "modal"
          }
        >
          <div class="modal-background" />
          <div class="modal-card">
            <header class="modal-card-head">
              <p class="modal-card-title">修改文章</p>
              <button
                class="delete"
                aria-label="close"
                onClick={event => this.closeUpdatePostModal()}
              />
            </header>
            <section class="modal-card-body">
              <PostEditor
                value={
                  this.state.postOnUpdate ? this.state.postOnUpdate.content : ""
                }
                changed={value => this.postOnUpdateChangeHandler(value)}
              />
            </section>
            <footer class="modal-card-foot">
              <button
                class="button is-primary"
                onClick={event => this.updatePostHandler(postID)}
              >
                確定
              </button>
              <button
                class="button"
                onClick={event => this.closeUpdatePostModal()}
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
