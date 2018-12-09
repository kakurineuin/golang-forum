import React, { Component } from "react";
import { Link } from "react-router-dom";
import dateFns from "date-fns";
import axios from "axios";
import produce from "immer";
import ForumStatistics from "../components/ForumStatistics";

/**
  首頁。
*/
class Home extends Component {
  state = {
    golang: {
      totalCount: 0,
      replyCount: 0,
      lastPostUsername: null,
      lastPostTime: null
    },
    nodeJS: {
      totalCount: 0,
      replyCount: 0,
      lastPostUsername: null,
      lastPostTime: null
    }
  };

  componentDidMount() {
    axios.get("/api/topics/statistics").then(response => {
      this.setState(
        produce(draft => {
          draft.golang = response.data.golang;
          draft.nodeJS = response.data.nodeJS;
        })
      );
    });
  }

  render() {
    const golang = this.state.golang;
    let golangLastPost = null;

    if (golang.lastPostTime) {
      golangLastPost = (
        <div>
          {dateFns.format(new Date(golang.lastPostTime), "YYYY/MM/DD HH:mm:ss")}
          <br />
          {golang.lastPostUsername}
        </div>
      );
    }

    const nodeJS = this.state.nodeJS;
    let nodeJSLastPost = null;

    if (nodeJS.lastPostTime) {
      nodeJSLastPost = (
        <div>
          {dateFns.format(new Date(nodeJS.lastPostTime), "YYYY/MM/DD HH:mm:ss")}
          <br />
          {nodeJS.lastPostUsername}
        </div>
      );
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
                <th>最新發文</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>
                  <Link
                    to={{
                      pathname: "/topics/golang",
                      state: { title: "Golang" }
                    }}
                  >
                    Golang
                  </Link>
                </td>
                <td>{golang.topicCount}</td>
                <td>{golang.replyCount}</td>
                <td>{golangLastPost}</td>
              </tr>
              <tr>
                <td>
                  <Link
                    to={{
                      pathname: "/topics/nodejs",
                      state: { title: "Node.js" }
                    }}
                  >
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
          <ForumStatistics />
        </div>
      </div>
    );
  }
}

export default Home;
