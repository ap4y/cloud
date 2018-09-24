import React, { Component } from 'react';

class FileBrowser extends Component {
  constructor(props) {
    super(props);
    this.state = { directories: [] };
  }

  componentDidMount() {
    this.setState({
      directories: [{ name: 'Test Directory' }, { name: 'Test Directory 2' }]
    });
  }

  render() {
    return (
      <div>
        <p>Directory</p>
      </div>
    );
  }
}

export default FileBrowser;
