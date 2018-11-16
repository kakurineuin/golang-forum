import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import * as postActions from '../modules/post';
import dateFns from 'date-fns';

class Home extends Component {
  componentDidMount() {
    this.props.onFindPostsStatistics();
  }

  render() {
    const golang = this.props.postsStatistics.golang;
    let golangLastPost = '';

    if (golang.lastPostTime) {
      golangLastPost = golang.lastPostAccount 
        + ' 於 ' 
        + dateFns.format(new Date(golang.lastPostTime), 'YYYY/MM/DD HH:mm:ss');
    }

    const nodeJS = this.props.postsStatistics.nodeJS;
    let nodeJSLastPost = '';

    if (nodeJS.lastPostTime) {
      nodeJSLastPost = nodeJS.lastPostAccount 
        + ' 於 ' 
        + dateFns.format(new Date(nodeJS.lastPostTime), 'YYYY/MM/DD HH:mm:ss');
    }

    return (
      <div className="columns">
        <div className="column is-9">
          <table className="table is-bordered is-striped is-hoverable is-fullwidth">
            <thead>
              <tr>
                <th>分類</th>
                <th>主題數</th>
                <th>回覆數</th>
                <th>最後新增</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>
                  <Link to="/posts/golang">
                    Golang
                  </Link>
                </td>
                <td>{golang.topicCount}</td>
                <td>{golang.replyCount}</td>
                <td>{golangLastPost}</td>
              </tr>
              <tr>
                <td>
                  <Link to="/posts/nodejs">
                    Node.js
                  </Link>
                </td>
                <td>{nodeJS.topicCount}</td>
                <td>{nodeJS.replyCount}</td>
                <td>{nodeJSLastPost}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div className="column is-3">
          <nav className="panel">
            <p className="panel-heading">
              TODO: 待實作。
            </p>
            <div className="panel-block">
            </div>
          </nav>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    postsStatistics: state.post.postsStatistics
  }
}

const mapDispatchToProps = dispatch => {
  return {
    onFindPostsStatistics: () => dispatch(postActions.findPostsStatistics())
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Home);