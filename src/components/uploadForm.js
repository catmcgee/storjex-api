import React from 'react';

class UploadForm extends React.Component {
  constructor() {
    super();
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleSubmit(event) {
    event.preventDefault();
    const data = new FormData(event.target);
    console.log(data)
    
    fetch('http://localhost:10000/api/v1', {
      method: 'POST',
      body: data,
    }).then(response => response.json())
    .then(data => alert("Your passphrase is " + data.password));
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit} encType = "multipart/form-data" >
        <label htmlFor="name">Give it a name</label>
        <input id="name" name="name" type="text"></input>
        <label htmlFor="file">Enter file</label>
        <input id="file" name="file" type="file" />

        <button>Upload file!</button>
      </form>
    );
  }
}

export default UploadForm;