import produce from "immer";
import _ from 'lodash';

// Action Types
export const SHOW = 'forum/message/SHOW';
export const DELETE = 'forum/message/DELETE';

// Reducer
export default function reducer(state = { list:[] }, action = {}) {
  return produce(state, draft => {
    switch (action.type) {
      case SHOW:
        draft.list.push(action.message);
        break;
      case DELETE:
        _.remove(draft.list, (message) => {
          return message.id === action.id;
        }) 
        break;
      default:
        break;
    }
  });
}

// Action Creators
export function showMessage(id, isError, text) {
  return {
    type: SHOW,
    message: {
      id,
      isError,
      text
    }
  };
};

export function deleteMessage(id) {
  return {
    type: DELETE,
    id
  };
};