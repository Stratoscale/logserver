import React, {Component} from 'react'
import {List} from 'antd'
import {withLoader} from 'utils'
import {send} from 'sockets/socket-actions'
import {connect} from 'react-redux'
import {API_ACTIONS} from 'consts'
import {createStructuredSelector} from 'reselect'
import {filesSelector} from 'selectors'

@connect(createStructuredSelector({
  files: filesSelector,
}), {
  send,
})
class FileTree extends Component {
  componentDidMount() {
    this.props.send(API_ACTIONS.GET_FILE_TREE, {
      base_path: [],
    })
  }

  render() {
    const {files} = this.props

    return (
      <List
        dataSource={files.toJS()}
        renderItem={item => <List.Item>{item.name}</List.Item>}
      />
    )
  }
}

export default withLoader(FileTree)