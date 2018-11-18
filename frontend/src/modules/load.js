import produce from "immer";

// Action Types
export const START = "forum/load/START";
export const STOP = "forum/load/STOP";

// Reducer
export default function reducer(state = {}, action = {}) {
  return produce(state, draft => {
    switch (action.type) {
      case START:
        draft.loading = true;
        break;
      case STOP:
        draft.loading = false;
        break;
      default:
        break;
    }
  });
}

// Action Creators
export function startLoad() {
  return {
    type: START
  };
}

export function stopLoad() {
  return {
    type: STOP
  };
}
