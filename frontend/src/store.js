import { createStore, applyMiddleware, compose, combineReducers } from 'redux';
import thunk from 'redux-thunk';
import auth from './modules/auth';
import load from './modules/load';
import message from './modules/message';
import post from './modules/post';

// 啟用 redux devTools chrome 擴充套件。
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

const rootReducer = combineReducers({
  auth,
  load,
  message,
  post
});

const store = createStore(rootReducer, composeEnhancers(
  applyMiddleware(thunk)
));

export default store;