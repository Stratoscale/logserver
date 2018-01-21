import React, {Component} from 'react'
import {List} from 'immutable'
import {contentSelector, filesSelector, hasPendingRequest, indexSelector, locationSelect} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import queryString from 'query-string'
import {clearContent, send, setSearch} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'
import {LinesView} from 'file-view/lines-view'
import {FSBar} from 'fs-bar'
import {navigate} from 'utils'
import Loader from 'loader/loader'

@connect(createStructuredSelector({
  location:   locationSelect,
  content:    contentSelector,
  files:      filesSelector,
  index:      indexSelector,
  requesting: hasPendingRequest,
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
    const {location, index, path, send} = this.props
    const search                        = queryString.parse(location.search)
    const {fs}                          = search

    if (fs) {
      this.setState({
        activeFs: [fs],
      })
    } else {
      const instance = index.getIn([path.join('/'), 'instances'], List()).first()

      this.setState({
        activeFs: [instance.get('fs')],
      })
    }

    this.props.setSearch('')
    this.props.clearContent()

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

  _handleToggle = (index) => {
    const {location: {pathname}} = this.props

    const fs = queryString.stringify({
      fs: index,
    })

    const activeFs = [index]

    this.setState({
      activeFs,
    })

    navigate(`${pathname}?${fs}`)
  }

  render() {
    const {content, location, index, requesting} = this.props
    let contentComponent                         = null
    if (requesting) {
      contentComponent = <Loader/>
    } else if (content && content.size > 0) {
      contentComponent = <LinesView lines={content} location={location}/>
    } else {
      contentComponent = <div>File is empty</div>
    }

    const path = location.pathname.split('/').filter(Boolean)
    const file = index.get(path.join('/'))

    return (
      <div className="file-view">
        <FSBar
          items={file.get('instances', List()).map(instance => ({
            name:   instance.get('fs'),
            active: this.state.activeFs.includes(instance.get('fs')),
          })).toJS()}
          onToggle={this._handleToggle}
        />
        {contentComponent}
      </div>
    )
  }
}

export default FileView