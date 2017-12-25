import React, {Component} from 'react'
import {Map} from 'immutable'
import {List, Icon} from 'antd'
import {withLoader} from 'utils'
import {send} from 'sockets/socket-actions'
import {connect} from 'react-redux'
import {API_ACTIONS} from 'consts'
import {createStructuredSelector} from 'reselect'
import {filesSelector} from 'selectors'

const File = ({path, isDir}) => {
  return <List.Item className="file"><Icon type={isDir ? 'folder' : 'file'}/> {path.join('/')}</List.Item>
}

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
        dataSource={files.get('files', Map()).valueSeq().toJS()}
        renderItem={File}
      />
    )
  }
}

export default withLoader(FileTree)