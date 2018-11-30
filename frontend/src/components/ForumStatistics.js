import React, { Component } from "react";
import axios from "axios";
import produce from "immer";

/**
 * 論壇統計資料。
 */
class ForumStatistics extends Component {
  state = {
    topicCount: 0, // 主題總數。
    replyCount: 0, // 回覆總數。
    userCount: 0 // 使用者總數。
  };

  componentDidMount() {
    axios.get("/api/forum/statistics").then(response => {
      const forumStatistics = response.data.forumStatistics;
      this.setState(
        produce(draft => {
          draft.topicCount = forumStatistics.topicCount;
          draft.replyCount = forumStatistics.replyCount;
          draft.userCount = forumStatistics.userCount;
        })
      );
    });
  }

  render() {
    return (
      <div className="card">
        <header className="card-header">
          <p className="card-header-title">論壇統計資料</p>
        </header>
        <div className="card-content">
          <div className="level">
            <div className="level-item has-text-centered">
              <div>
                <p className="heading">主題總數</p>
                <p className="title">{this.state.topicCount}</p>
              </div>
            </div>
            <div className="level-item has-text-centered">
              <div>
                <p className="heading">回覆總數</p>
                <p className="title">{this.state.replyCount}</p>
              </div>
            </div>
            <div className="level-item has-text-centered">
              <div>
                <p className="heading">使用者總數</p>
                <p className="title">{this.state.userCount}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default ForumStatistics;
