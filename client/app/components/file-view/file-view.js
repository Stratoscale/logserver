import React, {Component} from 'react'
import {List} from 'immutable'
import _uniq from 'lodash/uniq'
import {contentSelector, filesSelector, indexSelector, locationSelect} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import queryString from 'query-string'
import {clearContent, send, setSearch} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'
import LinesView from 'file-view/lines-view'
import {FSBar} from 'fs-bar'

@connect(createStructuredSelector({
  location: locationSelect,
  content:  contentSelector,
  files:    filesSelector,
  index:    indexSelector,
}), {
  send,
  setSearch,
  clearContent,
})
class FileView extends Component {
  constructor(props) {
    super(props)
    this.state = {
      activeFs: [],
    }
  }

  componentWillMount() {
    const {location} = this.props
    const search     = queryString.parse(location.search)
    const {fs}       = search
    this.setState({
      activeFs: [fs],
    })

    this.props.setSearch('')
    this.props.clearContent()
  }

  componentDidMount() {
    const {send, path} = this.props

    send(API_ACTIONS.GET_FILE_TREE, {
      base_path: [],
    })

    send(API_ACTIONS.GET_CONTENT, {
      path,
      filter_fs: this.state.activeFs,
    })
  }

  componentDidUpdate(prevProps, prevState) {
    if (prevState.activeFs !== this.state.activeFs) {
      const {send, path} = this.props
      this.props.clearContent()

      send(API_ACTIONS.GET_CONTENT, {
        path,
        filter_fs: this.state.activeFs,
      })
    }
  }

  _handleToggle = (index, value) => {
    let activeFs = [...this.state.activeFs]
    if (value) {
      activeFs.push(index)
      activeFs = _uniq(activeFs)
    } else {
      activeFs = activeFs.filter(v => v !== index)
    }

    this.setState({
      activeFs,
    })
  }

  render() {
    const {content, location, index} = this.props
    if (!content || content.size === 0) {
      return (
        <div>File is empty</div>
      )
    }

    const [_, ...path] = location.pathname.split('/').filter(Boolean)
    const file         = index.get(path.join('/'))

    return (
      <div>
        <FSBar
          items={file.get('instances', List()).map(instance => ({
            name:   instance.get('fs'),
            active: this.state.activeFs.includes(instance.get('fs')),
          }))}
          onToggle={this._handleToggle}
        />
        <LinesView lines={content}/>
      </div>
    )
  }
}

export default FileView