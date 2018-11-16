import React, { Component } from 'react';
import PostEditor from '../components/PostEditor';
import { connect } from 'react-redux';
import axios from 'axios';

class Posts extends Component {
  state = {
    title: '',
    content: ''
  };

  titleChangeHandler(event) {
    this.setState({ title: event.target.value });
  }

  contentChangeHandler(value) {
    this.setState({ content: value });
  }

  createTopicHandler() {
    console.log('props', this.props);
    console.log('state', this.state);
    axios.post(`/api/posts/${this.props.match.params.category}`, {
      userProfileID: this.props.user.id,
      title: this.state.title,
      content: this.state.content
    })
      .then(response => {
        console.log('create topic response', response);
        this.setState({
          title: '',
          content: ''
        });
      });
  }

  render() {
    let createTopic = null;

    if (this.props.user) {
      createTopic = (
        <div className="box">
          <div className="field">
            <label className="label">標題</label>
            <div className="control">
              <input type="text"
                className="input"
                placeholder="請輸入標題"
                value={this.state.title}
                onChange={event => this.titleChangeHandler(event)} />
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
                onClick={event => this.createTopicHandler()}>新增文章</button>
            </div>
          </div>
        </div>
      );
    }

    return (
      <div className="container">
        <table className="table is-bordered is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>Topic</th>
              <th>Created</th>
              <th>Statistics</th>
              <th>Last post</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td></td>
              <td></td>
              <td></td>
              <td></td>
            </tr>
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