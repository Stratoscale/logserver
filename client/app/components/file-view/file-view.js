import React, {Component} from 'react'
import {List, Set} from 'immutable'
import {contentSelector, filesSelector, hasPendingRequest, indexSelector, locationSelect} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import queryString from 'query-string'
import {clearContent, send, setSearch} from 'sockets/socket-actions'
import {API_ACTIONS, colorByLevel} from 'consts'
import {LinesView} from 'file-view/lines-view'
import {FSBar} from 'fs-bar'
import {navigate} from 'utils'
import Loader from 'loader/loader'
import {Checkbox} from 'antd'

const ALL_LEVELS = Set(['debug', 'info', 'warning', 'error'])

@connect(createStructuredSelector({
  location:   locationSelect,
  content:    contentSelector,
  files:      filesSelector,
  index:      indexSelector,
  requesting: hasPendingRequest(API_ACTIONS.GET_CONTENT),
}), {
  send,
  setSearch,
  clearContent,
})
class FileView extends Component {
  state = {
    activeFs:        [],
    showLevels:      ALL_LEVELS,
    showTimestamp:   true,
    showLinenumbers: true,
  }

  componentWillMount() {
    const {location, index, path, send} = this.props
    const search                        = queryString.parse(location.search)
    const {fs}                          = search
    let activeFs                        = []

    if (fs) {
      activeFs = [fs]
    } else {
      const instance = index.getIn([path.join('/'), 'instances'], List()).first()
      activeFs       = [instance.get('fs')]
    }

    this.setState({
      activeFs,
    })

    this.props.setSearch('')
    this.props.clearContent()

    send(API_ACTIONS.GET_CONTENT, {
      path,
      filter_fs: activeFs,
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

  _handleLevelToggle = (index) => {
    if (this.state.showLevels.includes(index)) {
      this.setState({
        showLevels: this.state.showLevels.delete(index),
      })
    } else {
      this.setState({
        showLevels: this.state.showLevels.add(index),
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
      contentComponent =
        <LinesView lines={content}
                   location={location}
                   showLevels={this.state.showLevels}
                   showLinenumbers={this.state.showLinenumbers}
                   showTimestamp={this.state.showTimestamp}/>
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
        <div>
          <FSBar onToggle={this._handleLevelToggle}
                 className="levels"
                 items={ALL_LEVELS.map(level => ({
                   name:   level,
                   active: this.state.showLevels.includes(level),
                   color:  colorByLevel(level),
                 })).toJS()}/>
          <Checkbox checked={this.state.showLinenumbers} onChange={({target: {checked}}) => {
            this.setState({
              showLinenumbers: checked,
            })
          }}>Line Numbers</Checkbox>

          <Checkbox checked={this.state.showTimestamp} onChange={({target: {checked}}) => {
            this.setState({
              showTimestamp: checked,
            })
          }}>Timestamps</Checkbox>


        </div>
        {contentComponent}
      </div>
    )
  }
}

export default FileView