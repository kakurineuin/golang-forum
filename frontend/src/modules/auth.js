import produce from "immer";
import axios from 'axios';
import format from 'date-fns/format';

// Action Types
export const AUTH_SUCCESS = 'forum/auth/AUTH_SUCCESS';
export const AUTH_LOGOUT = 'forum/auth/AUTH_LOGOUT';

// Reducer
export default function reducer(state = {}, action = {}) {
  console.log('action in reducer', action);

  return produce(state, draft => {
    switch (action.type) {
      case AUTH_SUCCESS:
        draft.user = action.user;
        break;
      case AUTH_LOGOUT:
        delete draft.user;
        break;
      default:
        break;
    }
  });
}

// Action Creators
export function authSuccess(user) {
  localStorage.setItem('user', JSON.stringify(user));
  return {
    type: AUTH_SUCCESS,
    user
  };
};

export function register(account, email, password) {
  return dispatch => {
    axios.post('/api/auth/register', { account, email, password })
      .then(response => {
        dispatch(authSuccess(buildUser(response)));
      });
  };
};

export function login(email, password) {
  return dispatch => {
    axios.post('/api/auth/login', { email, password })
      .then(response => {
        dispatch(authSuccess(buildUser(response)));
      });
  };
};

export function logout() {
  localStorage.removeItem('user');
  // TODO: 改為送出請求到後端把 token 加入到黑名單。
  return {
    type: AUTH_LOGOUT
  };
};

export function authCheckState() {
  return dispatch => {
    const user = JSON.parse(localStorage.getItem('user'));
    if (!user || user.expDate <= new Date()) {
      dispatch(logout());
    } else {
      dispatch(authSuccess(user));
    }
  };
};

function buildUser(response) {
  const user = response.data.userProfile;
  user.expDate = new Date(response.data.exp * 1000);
  console.log('expDate', format(user.expDate, 'YYYY/MM/DD HH:mm:ss'))
  user.token = response.data.token;
  return user;
}