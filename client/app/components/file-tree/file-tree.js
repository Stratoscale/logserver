import React, {Component} from 'react'
import {Map} from 'immutable'
import {List, Icon, Tag} from 'antd'
import {send} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {filesSelector, filterSelector, locationSelect} from 'selectors'
import {setCurrentPath} from 'file-tree/file-actions'
import {Link} from 'react-router-dom'

const File = ({path, is_dir, instances}) => {
  const last = path[path.length - 1]
  let content
  if (is_dir) {
    content = <span><Icon type={'folder'}/><Link to={`/files/${path.join('/')}`}>{last}</Link></span>
  } else {

    const viewURL = `/view?file=/${path.join('/')}`
    content       = <span>
      <Icon type={'file'}/> <Link to={viewURL}>{last}</Link>
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
    const {files, match: {params}, filter} = this.props

    const path         = (params[0] || '').split('/').filter(Boolean)
    const data         = files.getIn(path.concat(['files']), Map()).valueSeq().toJS()
    const filteredData = filter ? data.filter(file => file.key.includes(filter)) : data

    return (
      <List
        dataSource={filteredData}
        renderItem={file => <File {...file}/>}
      />
    )
  }
}

export default FileTree