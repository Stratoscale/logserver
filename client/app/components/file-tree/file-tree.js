import React, {Component} from 'react'
import {Map} from 'immutable'
import {List, Icon, Tag} from 'antd'
import {send} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {filesSelector, filterSelector, hasPendingRequest, indexSelector, locationSelect} from 'selectors'
import {setCurrentPath} from 'file-tree/file-actions'
import {Link} from 'react-router-dom'
import FileView from 'file-view/file-view'
import Loader from 'loader/loader'
import filesize from 'file-size'


const File = ({path, is_dir, instances, key, showFullPath = false}) => {
  const filename = showFullPath ? '/' + path.join('/') : path[path.length - 1]
  let content
  if (is_dir) {
    content = <span>
      <Icon type={'folder'}/><Link to={`/${path.join('/')}`}>{filename}</Link>
      {instances.map(instance => <Tag key={instance.fs}>{instance.fs}</Tag>)}
    </span>
  } else {
    const viewURL = `/${path.join('/')}`
    content       = <span>
      <Icon type={'file'}/> <Link to={viewURL}>{filename}</Link>
      {instances.map(instance => <Tag key={instance.fs}><Link to={`${viewURL}?fs=${instance.fs}`}>{instance.fs} <span className="size">({filesize(instance.size, {fixed: 0}).human()})</span></Link></Tag>)}
      <a href={`${window.location.origin}${window.__INIT__.basePath}/_dl/${path.join('/')}`} target="_blank" title="Download logs"><Icon type="download"/></a>
    </span>
  }
  return (
    <List.Item className="file">
      {content}
    </List.Item>
  )
}

@connect(createStructuredSelector({
  files:      filesSelector,
  index:      indexSelector,
  location:   locationSelect,
  filter:     filterSelector,
  requesting: hasPendingRequest(API_ACTIONS.GET_FILE_TREE),
}), {
  send,
  setCurrentPath,
})
class FileTree extends Component {
  render() {
    const {index, files, requesting, match: {params}, filter} = this.props

    const path  = (params[0] || '').split('/').filter(Boolean)
    const isDir = path.length === 0 || index.getIn([path.join('/'), 'is_dir'])
    let data
    let content = <FileView path={path} key={path}/>
    if (isDir) {
      data = files.getIn(path.concat(['files']), Map()).valueSeq().toJS()

      if (filter) {
        if (path.length) {
          data = index.filter((value, key) => key.startsWith(path.join('/') + '/'))
        } else {
          data = index
        }
        data = data.filter((value, key) => key.includes(filter)).valueSeq().toJS()
      }
      if (requesting) {
        content = <Loader/>
      } else {
        content = <List
          dataSource={data}
          renderItem={file => <File {...file} showFullPath={filter}/>}
        />
      }
    }

    return content
  }
}

export default FileTree