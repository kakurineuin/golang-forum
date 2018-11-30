import React, { Component } from "react";
import { Link } from "react-router-dom";
import PostEditor from "../components/PostEditor";
import Pagination from "../components/Pagination";
import { connect } from "react-redux";
import axios from "axios";
import produce from "immer";
import dateFns from "date-fns";

/**
  主題列表頁面。
*/
class Topics extends Component {
  state = {
    topic: "", // 新增文章的主題。
    content: "", // 新增文章的內文。
    topics: [], // 主題列表。
    totalCount: 0, // 主題總筆數。
    inputSearchTopic: "", // 搜尋輸入框輸入的值。
    searchTopic: "", // 請求參數：搜尋的主題。
    paginationKey: Math.random() // 用來觸發分頁重新 render 並查詢資料。
  };

  topicChangeHandler(event) {
    const value = event.target.value;
    this.setState(
      produce(draft => {
        draft.topic = value;
      })
    );
  }

  contentChangeHandler(value) {
    this.setState(
      produce(draft => {
        draft.content = value;
      })
    );
  }

  createTopicHandler() {
    axios
      .post(`/api/topics/${this.props.match.params.category}`, {
        userProfileID: this.props.user.id,
        topic: this.state.topic,
        content: this.state.content
      })
      .then(response => {
        this.setState(
          produce(draft => {
            draft.topic = "";
            draft.content = "";
            draft.inputSearchTopic = "";
            draft.searchTopic = "";
            draft.paginationKey = Math.random(); // 觸發分頁重新查詢。
          })
        );
      });
  }

  inputSearchTopicChangeHandler(value) {
    this.setState(
      produce(draft => {
        draft.inputSearchTopic = value;
      })
    );
  }

  inputSearchTopicKeyPressHandler(event) {
    if (event.key !== "Enter") return;
    this.searchHandler();
  }

  searchHandler() {
    this.setState(
      produce(draft => {
        draft.searchTopic = draft.inputSearchTopic;
        draft.inputSearchTopic = "";
        draft.paginationKey = Math.random(); // 觸發分頁重新查詢。
      })
    );
  }

  deleteTagHandler() {
    this.setState(
      produce(draft => {
        draft.searchTopic = "";
        draft.paginationKey = Math.random(); // 觸發分頁重新查詢。
      })
    );
  }

  findPostsTopics(offset, limit) {
    axios
      .get(`/api/topics/${this.props.match.params.category}`, {
        params: {
          searchTopic: this.state.searchTopic,
          offset,
          limit
        }
      })
      .then(response => {
        this.setState(
          produce(draft => {
            draft.topics = response.data.topics;
            draft.totalCount = response.data.totalCount;
          })
        );
      });
  }

  render() {
    let createTopic = null;

    if (this.props.user) {
      createTopic = (
        <div className="box">
          <div className="field">
            <label className="label">主題</label>
            <div className="control">
              <input
                type="text"
                className="input"
                placeholder="請輸入主題"
                value={this.state.topic}
                onChange={event => this.topicChangeHandler(event)}
              />
            </div>
          </div>
          <div className="field">
            <label className="label">內文</label>
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
                onClick={event => this.createTopicHandler()}
              >
                新增主題
              </button>
            </div>
          </div>
        </div>
      );
    }

    const topics = this.state.topics.map((post, index) => {
      let lastReply = <td />;

      if (post.lastReplyUsername) {
        lastReply = (
          <td>
            {dateFns.format(
              new Date(post.lastReplyCreatedAt),
              "YYYY/MM/DD HH:mm:ss"
            )}
            <br />
            {post.lastReplyUsername}
          </td>
        );
      }

      const category = this.props.match.params.category;

      return (
        <tr key={post.id}>
          <td>
            <Link to={`/topics/${category}/${post.id}`}>{post.topic}</Link>
          </td>
          <td>{post.replyCount}</td>
          <td>
            {dateFns.format(new Date(post.createdAt), "YYYY/MM/DD HH:mm:ss")}
            <br />
            {post.username}
          </td>
          {lastReply}
        </tr>
      );
    });

    // 顯示搜尋主題的標籤。
    let tagSearchTopic = null;

    if (this.state.searchTopic) {
      tagSearchTopic = (
        <span className="tag is-medium">
          {this.state.searchTopic}
          <button
            className="delete is-small"
            onClick={event => this.deleteTagHandler()}
          />
        </span>
      );
    }

    return (
      <div className="container">
        <div className="field has-addons">
          <div className="control">
            <input
              className="input"
              type="text"
              value={this.state.inputSearchTopic}
              onChange={event =>
                this.inputSearchTopicChangeHandler(event.target.value)
              }
              onKeyPress={event => this.inputSearchTopicKeyPressHandler(event)}
              placeholder="搜尋主題"
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
        {tagSearchTopic}
        <hr />
        <table className="table is-bordered is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>主題</th>
              <th>回覆數</th>
              <th>作者</th>
              <th>最新回覆</th>
            </tr>
          </thead>
          <tbody>{topics}</tbody>
        </table>
        <Pagination
          key={this.state.paginationKey}
          totalCount={this.state.totalCount}
          findData={(offset, limit) => this.findPostsTopics(offset, limit)}
        />
        <br />
        {createTopic}
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
)(Topics);
