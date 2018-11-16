import produce from "immer";
import _ from 'lodash';
import axios from 'axios';

// Action Types
export const FIND_POSTS_STATISTICS_SUCCESS = 'forum/post/FIND_POSTS_STATISTICS_SUCCESS';

// Reducer
export default function reducer(state = {
  postsStatistics: {
    golang: {
      totalCount: 0,
      replyCount: 0,
      lastPostAccount: null,
      lastPostTime: null
    },
    nodeJS: {
      totalCount: 0,
      replyCount: 0,
      lastPostAccount: null,
      lastPostTime: null
    },
  }
}, action = {}) {
  return produce(state, draft => {
    switch (action.type) {
      case FIND_POSTS_STATISTICS_SUCCESS:
        draft.postsStatistics = action.postsStatistics
        break;
      default:
        break;
    }
  });
}

// Action Creators
export function findPostsStatisticsSuccess(postsStatistics) {
  return {
    type: FIND_POSTS_STATISTICS_SUCCESS,
    postsStatistics
  }
}

export function findPostsStatistics() {
  return dispatch => {
    axios.get('/api/posts/statistics')
      .then(response => {
        dispatch(findPostsStatisticsSuccess(response.data));
      });
  };
};
