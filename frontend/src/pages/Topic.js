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
    content: "", // 回覆內文。
    paginationKey: Math.random() // 用來觸發分頁重新 render 並查詢資料。
  };

  contentChangeHandler(value) {
    this.setState(
      produce(draft => {
        draft.content = value;
      })
    );
  }

  createReplyHandler() {
    console.log("props", this.props);
    console.log("state", this.state);
    axios
      .post(`/api/posts/${this.props.match.params.category}`, {
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

  findPostsByTopicID(offset, limit) {
    const category = this.props.match.params.category;
    const id = this.props.match.params.id;
    axios
      .get(`/api/posts/${category}/topics/${id}`, {
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

    const id = parseInt(this.props.match.params.id, 10);
    const posts = this.state.posts.map((post, index) => {
      let content = null;

      // 若是主題文章，顯示主題。
      content = (
        <td>
          {post.id === id ? (
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
            <h6 className="title is-6">{post.account}</h6>
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
