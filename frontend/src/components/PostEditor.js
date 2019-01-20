import React, { Component } from "react";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";

/**
  文章編輯器元件。
*/
class PostEditor extends Component {
  constructor(props) {
    super(props);
    this.quillRef = React.createRef();
    this.modules = {
      toolbar: {
        container: [
          [{ header: [1, 2, false] }],
          ["bold", "italic", "underline", "strike", "blockquote", "code-block"],
          [
            { list: "ordered" },
            { list: "bullet" },
            { indent: "-1" },
            { indent: "+1" }
          ],
          ["image"],
          ["clean"]
        ],
        handlers: {
          image: (image, callback) => {
            var range = this.quillRef.current.getEditor().getSelection();
            var value = prompt("請輸入圖片網址。");
            if (value) {
              this.quillRef.current
                .getEditor()
                .insertEmbed(range.index, "image", value, "user");
            }
          }
        }
      }
    };
    this.formats = [
      "header",
      "bold",
      "italic",
      "underline",
      "strike",
      "blockquote",
      "code-block",
      "list",
      "bullet",
      "indent",
      "image"
    ];
  }

  changeHandler(value) {
    this.props.changed(value);
  }

  render() {
    return (
      <div style={{ marginBottom: "10px" }}>
        <ReactQuill
          {...this.props}
          ref={this.quillRef}
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
