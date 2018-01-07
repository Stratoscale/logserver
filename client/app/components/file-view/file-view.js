import React, {Component} from 'react'
import {Map, List} from 'immutable'
import {Tag} from 'antd'
import {contentSelector, filesSelector, locationSelect} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import queryString from 'query-string'
import {send, setSearch} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'
import LinesView from 'file-view/lines-view'

@connect(createStructuredSelector({
  location: locationSelect,
  content:  contentSelector,
  files:    filesSelector,
}), {
  send,
  setSearch,
})
class FileView extends Component {
  componentDidMount() {
    const {location, send} = this.props
    const search           = queryString.parse(location.search)
    const {fs, file = ''}  = search

    const filter_fs = fs ? [fs] : []
    this.props.setSearch('')
    send(API_ACTIONS.GET_FILE_TREE, {
      base_path: [],
    })
    send(API_ACTIONS.GET_CONTENT, {
      path: file.split('/').filter(Boolean),
      filter_fs,
    })
  }

  render() {
    const {content, location, files} = this.props
    if (!content) {
      return (
        <div>File is empty</div>
      )
    }

    const search = queryString.parse(location.search)
    const {searchFile = ''} = search
    const path   = searchFile.split('/').filter(Boolean)

    const file   = files.getIn(path.slice(0, -1).concat(['files', searchFile.slice(1)]), Map())

    return (
      <div>
        <div>
          {file.get('instances', List()).map(instance => <Tag>{instance.get('fs')}</Tag>)}
        </div>
        <LinesView lines={content}/>
      </div>
    )
  }
}

export default FileView