import React, {Component} from 'react'
import {Map} from 'immutable'
import {List, Icon, Tag} from 'antd'
import {withLoader} from 'utils'
import {send} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {filesSelector, locationSelect} from 'selectors'
import {setCurrentPath} from 'file-tree/file-actions'
import {Link} from 'react-router-dom'

const File = ({path, is_dir, fs}) => {
  const last = path[path.length - 1]
  return (
    <List.Item className="file">
      <Icon type={is_dir ? 'folder' : 'file'}/> <Link to={is_dir ? `/files/${path.join('/')}` : `/view?file=/${path.join('/')}`}>{last}</Link>
      {fs.map(node => <Tag key={node}>{node}</Tag>)}

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