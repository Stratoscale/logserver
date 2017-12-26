import React, {Component} from 'react'
import {Map, List} from 'immutable'
import {Tag} from 'antd'
import {contentSelector, filesSelector, locationSelect} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import queryString from 'query-string'
import {send} from 'sockets/socket-actions'
import cn from 'classnames'
import {API_ACTIONS} from 'consts'

@connect(createStructuredSelector({
  location: locationSelect,
  content:  contentSelector,
  files:    filesSelector,
}), {
  send,
})
class FileView extends Component {
  componentDidMount() {
    const {location, send} = this.props
    const search           = queryString.parse(location.search)

    send(API_ACTIONS.GET_FILE_TREE, {
      base_path: [],
    })
    send(API_ACTIONS.GET_CONTENT, {
      path: search.file.split('/').filter(Boolean),
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
        <div>{content.map((line, index) => <div key={index} className={cn(line.level)}>{line.msg}</div>)}</div>
      </div>
    )
  }
}

export default FileView