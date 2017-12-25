import React, {Component} from 'react'
import {withLoader} from 'utils'
import {send} from 'sockets/socket-actions'
import {connect} from 'react-redux'

@connect(null, {
  send
})
class FolderTree extends Component {
  componentDidMount() {
    this.props.send('get-file-tree')
  }

  render() {
    return (
      <h1>Folder Tree</h1>
    )
  }
}

export default withLoader(FolderTree)