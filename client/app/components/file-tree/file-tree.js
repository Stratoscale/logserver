import React, {Component} from 'react'
import {Map} from 'immutable'
import {List, Icon} from 'antd'
import {withLoader} from 'utils'
import {send} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {filesSelector, locationSelect} from 'selectors'
import {setCurrentPath} from 'file-tree/file-actions'
import {Link} from 'react-router-dom'

const File = ({path, is_dir}) => {
  const last = path[path.length - 1]
  return (
    <List.Item className="file">
      <Icon type={is_dir ? 'folder' : 'file'}/> <Link to={`/files/${path.join('/')}`}>{last}</Link>
    </List.Item>
  )
}

@connect(createStructuredSelector({
  files:    filesSelector,
  location: locationSelect,
}), {
  send,
  setCurrentPath,
})
class FileTree extends Component {
  componentDidMount() {
    this.props.send(API_ACTIONS.GET_FILE_TREE, {
      base_path: [],
    })
  }

  render() {
    const {files, match: {params}} = this.props

    const path = (params[0] || '').split('/').filter(Boolean)

    return (
      <List
        dataSource={files.getIn(path.concat(['files']), Map()).valueSeq().toJS()}
        renderItem={file => <File {...file}/>}
      />
    )
  }
}

export default withLoader(FileTree)