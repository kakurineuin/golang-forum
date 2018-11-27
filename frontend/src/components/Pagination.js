import React, { Component } from "react";
import produce from "immer";

/**
  分頁元件。
*/
class Pagination extends Component {
  state = {
    currentPage: 1, // 當前頁數。
    limit: 10 // 每頁筆數。
  };

  componentDidMount() {
    this.props.findData(this.getOffset(), this.state.limit);
  }

  getOffset() {
    return (this.state.currentPage - 1) * this.state.limit;
  }

  pageClickHandler(page) {
    this.setState(
      produce(draft => {
        draft.currentPage = page;
      }),
      () => {
        this.props.findData(this.getOffset(), this.state.limit);
      }
    );
  }

  render() {
    const currentPage = this.state.currentPage; // 當前頁數。
    const totalCount = this.props.totalCount; // 總筆數。
    const maxPage = Math.ceil(totalCount / this.state.limit) || 1; // 最大頁數。
    console.log(
      "currentPage: {}, totalCount: {}, maxPage: {}",
      currentPage,
      totalCount,
      maxPage
    );

    const pages = [];

    // 第一頁。
    if (currentPage === 1) {
      pages.push(
        <li key={1}>
          <a
            className="pagination-link is-current"
            aria-label="Page 1"
            aria-current="page"
            onClick={() => this.pageClickHandler(1)}
          >
            1
          </a>
        </li>
      );
    } else {
      pages.push(
        <li key={1}>
          <a
            className="pagination-link"
            aria-label="Goto page 1"
            onClick={() => this.pageClickHandler(1)}
          >
            1
          </a>
        </li>
      );
    }

    // 顯示省略符號。
    if (currentPage - 3 > 1) {
      pages.push(
        <li key={"ellipsis1"}>
          <span className="pagination-ellipsis">&hellip;</span>
        </li>
      );
    }

    // 顯示當前頁的前兩頁和後兩頁。
    for (
      let i = Math.max(currentPage - 2, 2);
      i <= Math.min(currentPage + 2, maxPage);
      i++
    ) {
      if (i === currentPage) {
        pages.push(
          <li key={i}>
            <a
              className="pagination-link is-current"
              aria-label={`Page ${i}`}
              aria-current="page"
              onClick={() => this.pageClickHandler(i)}
            >
              {i}
            </a>
          </li>
        );
        continue;
      }

      pages.push(
        <li key={i}>
          <a
            className="pagination-link"
            aria-label={`Goto page ${i}`}
            onClick={() => this.pageClickHandler(i)}
          >
            {i}
          </a>
        </li>
      );
    }

    // 顯示省略符號。
    if (currentPage + 3 < maxPage) {
      pages.push(
        <li key={"ellipsis2"}>
          <span className="pagination-ellipsis">&hellip;</span>
        </li>
      );
    }

    // 最後一頁。
    if (currentPage + 2 < maxPage) {
      pages.push(
        <li key={maxPage}>
          <a
            className="pagination-link"
            aria-label={`Goto page ${maxPage}`}
            onClick={() => this.pageClickHandler(maxPage)}
          >
            {maxPage}
          </a>
        </li>
      );
    }

    return (
      <nav className="pagination" role="navigation" aria-label="pagination">
        <button
          className="button pagination-previous"
          disabled={currentPage === 1}
          onClick={() => this.pageClickHandler(currentPage - 1)}
        >
          上一頁
        </button>
        <button
          className="button pagination-next"
          disabled={currentPage === maxPage}
          onClick={() => this.pageClickHandler(currentPage + 1)}
        >
          下一頁
        </button>
        <ul className="pagination-list">{pages}</ul>
      </nav>
    );
  }
}

export default Pagination;
