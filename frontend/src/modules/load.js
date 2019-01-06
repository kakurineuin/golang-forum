import produce from "immer";

// Action Types
export const START = "forum/load/START";
export const STOP = "forum/load/STOP";

// Reducer
export default function reducer(
  state = { ajaxCount: 0, loading: false },
  action = {}
) {
  return produce(state, draft => {
    switch (action.type) {
      case START:
        draft.ajaxCount++;
        break;
      case STOP:
        draft.ajaxCount--;
        break;
      default:
        break;
    }
    draft.loading = draft.ajaxCount > 0;
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
