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
    topicID: null, // 主題文章 ID。
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
        replyPostID: this.state.posts[0].id
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

    const posts = this.state.posts.map((post, index) => {
      // TODO: 列出討論串文章。
    });

    return (
      <div className="container">
        <table className="table is-bordered is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>作者</th>
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
