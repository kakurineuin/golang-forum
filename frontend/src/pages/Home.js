import React, { Component } from 'react';
import { Link } from 'react-router-dom';

class Home extends Component {
  render() {
    return (
      <div className="columns">
        <div className="column is-9">
          <table className="table is-bordered is-striped is-hoverable is-fullwidth">
            <thead>
              <tr>
                <th>Category</th>
                <th>Topics</th>
                <th>Posts</th>
                <th>Last post</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>
                  <Link to="/posts/golang">
                    Golang
                  </Link>
                </td>
                <td></td>
                <td></td>
                <td></td>
              </tr>
              <tr>
                <td>
                  <Link to="/posts/nodejs">
                    Node.js
                  </Link>
                </td>
                <td></td>
                <td></td>
                <td></td>
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

export default Home;