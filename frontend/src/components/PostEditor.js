import React, { Component } from "react";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";

/**
  文章編輯器元件。
*/
class PostEditor extends Component {
  modules = {
    toolbar: [
      [{ header: [1, 2, false] }],
      ["bold", "italic", "underline", "strike", "blockquote", "code-block"],
      [
        { list: "ordered" },
        { list: "bullet" },
        { indent: "-1" },
        { indent: "+1" }
      ],
      ["clean"]
    ]
  };

  formats = [
    "header",
    "bold",
    "italic",
    "underline",
    "strike",
    "blockquote",
    "code-block",
    "list",
    "bullet",
    "indent"
  ];

  changeHandler(value) {
    this.props.changed(value);
  }

  render() {
    return (
      <div style={{ marginBottom: "10px" }}>
        <ReactQuill
          {...this.props}
          theme="snow"
          value={this.props.value}
          modules={this.modules}
          formats={this.formats}
          onChange={value => this.changeHandler(value)}
        />
      </div>
    );
  }
}

export default PostEditor;
