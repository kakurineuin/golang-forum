import React, { Component } from 'react';
import PostEditor from '../components/PostEditor';
import { connect } from 'react-redux';
import axios from 'axios';
import produce from "immer";
import dateFns from 'date-fns';

class Posts extends Component {
  state = {
    topic: '',
    content: '',
    offset: 0,
    limit: 10,
    posts: []
  };

  componentDidMount() {
    this.findPosts();
  }

  topicChangeHandler(event) {
    this.setState({ topic: event.target.value });
  }

  contentChangeHandler(value) {
    this.setState({ content: value });
  }

  createTopicHandler() {
    console.log('props', this.props);
    console.log('state', this.state);
    axios.post(`/api/posts/${this.props.match.params.category}`, {
      userProfileID: this.props.user.id,
      topic: this.state.topic,
      content: this.state.content
    })
      .then(response => {
        console.log('create topic response', response);
        this.setState(
          produce(draft => {
            draft.topic = '';
            draft.content = '';
          })
          , () => {
            this.findPosts();
          });
      });
  }

  findPosts() {
    axios.get(`/api/posts/${this.props.match.params.category}`, {
      params: {
        offset: this.state.offset,
        limit: this.state.limit
      }
    })
      .then(response => {
        this.setState(
          produce(draft => {
            draft.posts = response.data.posts
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
              <input type="text"
                className="input"
                placeholder="請輸入主題"
                value={this.state.topic}
                onChange={event => this.topicChangeHandler(event)} />
            </div>
          </div>
          <div className="field">
            <label className="label">內文</label>
            <div className="control">
              <PostEditor value={this.state.content} changed={value => this.contentChangeHandler(value)} />
            </div>
          </div>
          <div className="field">
            <div className="control">
              <button className="button is-primary"
                onClick={event => this.createTopicHandler()}>新增主題</button>
            </div>
          </div>
        </div>
      );
    }

    const posts = this.state.posts.map((post, index) => {
      let lastReply = (<td></td>);

      if (post.lastReplyAccount) {
        lastReply = (
          <td>
            {dateFns.format(new Date(post.lastReplyCreatedAt), 'YYYY/MM/DD HH:mm:ss')}
            <br />
            {post.lastReplyAccount}
          </td>
        );
      }

      return (
        <tr>
          <td>{post.topic}</td>
          <td>{post.replyCount}</td>
          <td>
            {dateFns.format(new Date(post.createdAt), 'YYYY/MM/DD HH:mm:ss')}
            <br />
            {post.account}
          </td>
          {lastReply}
        </tr>
      )
    });

    return (
      <div className="container">
        <table className="table is-bordered is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>主題</th>
              <th>回覆數</th>
              <th>作者</th>
              <th>最新發文</th>
            </tr>
          </thead>
          <tbody>
            {posts}
          </tbody>
        </table>
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

export default connect(mapStateToProps, null)(Posts);