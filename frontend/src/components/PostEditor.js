import React, { Component } from 'react';
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';

class PostEditor extends Component {
  modules = {
    toolbar: [
      [{ 'header': [1, 2, false] }],
      ['bold', 'italic', 'underline', 'strike', 'blockquote', 'code-block'],
      [{ 'list': 'ordered' }, { 'list': 'bullet' }, { 'indent': '-1' }, { 'indent': '+1' }],
      ['link'],
      ['clean']
    ]
  };

  changeHandler(value) {
    this.props.changed(value);
  }

  render() {
    return (
      <div style={{marginBottom: '10px'}}>
      <ReactQuill theme="snow"
        value={this.props.value}
        modules={this.modules}
        onChange={value => this.changeHandler(value)} />
      </div>
    );
  }
}

export default PostEditor;
