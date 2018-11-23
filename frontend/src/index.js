import React from "react";
import ReactDOM from "react-dom";
import { Router } from "react-router-dom";
import history from "./utils/history";
import "./index.css";
import App from "./App";
import * as serviceWorker from "./serviceWorker";
import store from "./store";
import { Provider } from "react-redux";
import axios from "axios";
import { startLoad, stopLoad } from "./modules/load";
import { showMessage } from "./modules/message";
import "bulma/css/bulma.css";
import * as authActions from "./modules/auth";

// 初始化 axios。
axios.interceptors.request.use(
  config => {
    // Do something before request is sent
    store.dispatch(startLoad());
    const user = JSON.parse(localStorage.getItem("user"));

    if (user) {
      config.headers = { Authorization: "Bearer " + user.token };
    }

    return config;
  },
  error => {
    // Do something with request error
    // 顯示錯誤訊息。
    store.dispatch(showMessage(Date.now(), true, error.toString()));
    store.dispatch(stopLoad());
    return Promise.reject(error);
  }
);

axios.interceptors.response.use(
  response => {
    // Do something with response data
    console.log("interceptor response", response);
    store.dispatch(stopLoad());

    // 顯示訊息。
    if (response.data && response.data.message) {
      store.dispatch(showMessage(Date.now(), false, response.data.message));
    }

    return response;
  },
  error => {
    // Do something with response error
    console.log("interceptor error.response", error.response);
    store.dispatch(stopLoad());

    // 若 JWT 過期。
    if (error.response.status === 401 && localStorage.getItem("user")) {
      store.dispatch(authActions.logout());
      store.dispatch(
        showMessage(Date.now(), true, "登入時效已過期，請重新登入。")
      );
    } else {
      // 顯示錯誤訊息。
      const text =
        error.response.data && error.response.data.message
          ? error.response.data.message
          : error.toString();
      store.dispatch(showMessage(Date.now(), true, text));
    }

    return Promise.reject(error);
  }
);

const app = (
  <Provider store={store}>
    <Router history={history}>
      <App />
    </Router>
  </Provider>
);

ReactDOM.render(app, document.getElementById("root"));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
