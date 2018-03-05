import React, {Component} from 'react'
import {List, Map} from 'immutable'
import {
  contentSelector, filesSelector, findSelector, hasPendingRequest, indexSelector, levelsSelector, locationSelect, matchesSelector,
} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import queryString from 'query-string'
import {clearContent, send, setLevels, setSearch} from 'sockets/socket-actions'
import {API_ACTIONS, colorByLevel} from 'consts'
import {LinesView} from 'file-view/lines-view'
import {FSBar} from 'fs-bar'
import {navigate} from 'utils'
import Loader from 'loader/loader'
import {Checkbox, Icon} from 'antd'
import filesize from 'file-size'
import {ALL_LEVELS} from 'reducers/app-reducers'

@connect(createStructuredSelector({
  location:   locationSelect,
  content:    contentSelector,
  matches:    matchesSelector,
  files:      filesSelector,
  levels:     levelsSelector,
  find:       findSelector,
  index:      indexSelector,
  requesting: hasPendingRequest(API_ACTIONS.GET_CONTENT),
}), {
  send,
  setSearch,
  clearContent,
  setLevels,
})
class FileView extends Component {
  state = {
    activeFs:        [],
    showTimestamp:   true,
    showLinenumbers: true,
    showThreadName:  false,
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
      if (instance) {
        activeFs = [instance.get('fs')]
      }
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
    const {levels, setLevels} = this.props
    if (levels.includes(index)) {
      setLevels(levels.delete(index))
    } else {
      setLevels(levels.add(index))
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

  _matchToLine = (matches, index) => {
    let matchCount = 0
    for (let lineIndex of matches.keySeq().sort()) {
      const lineMatches = matches.get(lineIndex)
      if (matchCount + lineMatches.length > index) {
        return lineIndex
      }
      matchCount += lineMatches.length
    }
  }

  render() {
    let {matches, content, location, index, requesting, find, levels} = this.props

    const scrollToLine = matches.size ? this._matchToLine(matches, find.get('index')) : undefined

    let contentComponent = null

    if (requesting) {
      contentComponent = <Loader/>
    } else if (content.size > 0) {
      contentComponent =
        <LinesView lines={content}
                   location={location}
                   matches={matches}
                   findIndex={find.get('index')}
                   findQuery={find.get('query')}
                   scrollToLine={scrollToLine}
                   showLevels={levels}
                   showLinenumbers={this.state.showLinenumbers}
                   showTimestamp={this.state.showTimestamp}
                   showThreadName={this.state.showThreadName}
        />
    } else {
      contentComponent = <div>File is empty</div>
    }

    const path = location.pathname.split('/').filter(Boolean)
    const file = index.get(path.join('/'), Map())

    return (
      <div className="file-view">
        <FSBar
          items={file.get('instances', List()).map(instance => ({
            name:    instance.get('fs'),
            content: <span>{instance.get('fs')} <span className="size">({filesize(instance.get('size'), {fixed: 0}).human()})</span></span>,
            active:  this.state.activeFs.includes(instance.get('fs')),
          })).toJS()}
          onToggle={this._handleToggle}
        />
        <div>
          <FSBar onToggle={this._handleLevelToggle}
                 className="levels"
                 items={ALL_LEVELS.map(level => ({
                   name:   level,
                   active: levels.includes(level),
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

          <Checkbox checked={this.state.showThreadName} onChange={({target: {checked}}) => {
            this.setState({
              showThreadName: checked,
            })
          }}>Thread Name</Checkbox>

          <a href={`${window.location.origin}${window.__INIT__.basePath}/_dl/${path.join('/')}?fs=${this.state.activeFs}`} target="_blank"><Icon type="eye"/> Show Original</a>


        </div>
        {contentComponent}
      </div>
    )
  }
}

export default FileView
