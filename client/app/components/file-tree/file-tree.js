import React, {Component} from 'react'
import {Map} from 'immutable'
import {List, Icon, Tag} from 'antd'
import {send} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {filesSelector, filterSelector, indexSelector, locationSelect} from 'selectors'
import {setCurrentPath} from 'file-tree/file-actions'
import {Link} from 'react-router-dom'

const File = ({path, is_dir, instances, key, showFullPath = false}) => {
  const filename = showFullPath ? '/' + path.join('/') : path[path.length - 1]
  let content
  if (is_dir) {
    content = <span>
      <Icon type={'folder'}/><Link to={`/files/${path.join('/')}`}>{filename}</Link>
      {instances.map(instance => <Tag key={instance.fs}>{instance.fs}</Tag>)}
    </span>
  } else {
    const viewURL = `/view?file=/${path.join('/')}`
    content       = <span>
      <Icon type={'file'}/> <Link to={viewURL}>{filename}</Link>
      {instances.map(instance => <Tag key={instance.fs}><Link to={`${viewURL}&fs=${instance.fs}`}>{instance.fs}</Link></Tag>)}
    </span>
  }
  return (
    <List.Item className="file">
      {content}
    </List.Item>
  )
}

@connect(createStructuredSelector({
  files:    filesSelector,
  index:    indexSelector,
  location: locationSelect,
  filter:   filterSelector,
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
    const {index, files, match: {params}, filter} = this.props

    const path = (params[0] || '').split('/').filter(Boolean)
    let data   = files.getIn(path.concat(['files']), Map()).valueSeq().toJS()
    if (filter) {
      if (path.length) {
        data = index.filter((value, key) => key.startsWith(path.join('/') + '/'))
      } else {
        data = index
      }
      data = data.filter((value, key) => key.includes(filter)).valueSeq().toJS()
    }

    return (
      <List
        dataSource={data}
        renderItem={file => <File {...file} showFullPath={filter}/>}
      />
    )
  }
}

export default FileTree