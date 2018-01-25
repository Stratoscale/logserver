import React, {Component} from 'react'
import _isString from 'lodash/isString'
import _flatMap from 'lodash/flatMap'
import cn from 'classnames'
import {Set, Map} from 'immutable'
import {Tag} from 'antd'
import PropTypes from 'prop-types'
import {Grid, AutoSizer} from 'react-virtualized'
import {Link} from 'react-router-dom'
import queryString from 'query-string'
import moment from 'moment'
import {colorByLevel} from 'consts'
import ImmutablePropTypes from 'react-immutable-proptypes'
import {navigate} from 'utils'

const calculateMaxLengths = (lines) => lines.reduce(({message, level}, line) => ({
  message: Math.max(message, Array.isArray(line.get('msg')) ? line.get('msg').reduce((lineLength, token) => {
    if (_isString(token)) {
      return lineLength + token.length
    } else {
      return lineLength + 1
    }
  }, 0) : line.get('msg').length),
  level:   Math.max(level, line.get('level').length),
}), {message: 0, level: 0})

const calculateLines = ({lines, showLevels, showFilename}) => {
  if (showFilename) {
    return lines.groupBy(line => line.get('file_name')).entrySeq().flatMap(([filename, groupLines]) => [Map({
      type: 'filename',
      msg:  filename,
      line: groupLines.getIn([0, 'line']),
      fs:   groupLines.getIn([0, 'fs']),
    }), ...groupLines]).toList()
  }
  return lines
}

class LinesView extends Component {
  state = {
    maxLengths: calculateMaxLengths(this.props.lines),
    lines:      calculateLines(this.props),
  }

  static propTypes = {
    showFilename:    PropTypes.bool,
    showTimestamp:   PropTypes.bool,
    showLinenumbers: PropTypes.bool,
    scrollToLine:    PropTypes.number,
    showLevels:      ImmutablePropTypes.set,
  }

  static defaultProps = {
    showFilename:    false,
    showTimestamp:   false,
    showLinenumbers: false,
    showLevels:      Set(),
  }

  componentWillReceiveProps({lines, showLevels, showFilename}) {
    if (!this.props.lines.equals(lines) || !this.props.showLevels.equals(showLevels)) {
      this.setState({
        maxLengths: calculateMaxLengths(lines),
        lines:      calculateLines({lines, showLevels, showFilename}),
      }, () => {
        if (this.grid) {
          this.grid.recomputeGridSize()
        }
      })
    } else {
      this.grid.recomputeGridSize()
    }
  }

  _getColumns = () => {
    const {maxLengths}                                   = this.state
    const {showFilename, showTimestamp, showLinenumbers} = this.props

    if (showFilename) {
      return [
        {
          width: maxLengths.message * 8,
          name:  'msg',
        },
      ]
    } else {
      const columns = [
        {
          width: maxLengths.level * 9,
          name:  'level',
        },
        {
          width: maxLengths.message * 8,
          name:  'msg',
        },
      ]

      if (showTimestamp) {
        columns.unshift({
          width: 125,
          name:  'timestamp',
        })
      }
      if (showLinenumbers) {
        columns.unshift({
          width: 40,
          name:  'linenumber',
        })
      }
      return columns
    }
  }

  _getColumnWidth = ({index}) => {
    return this._getColumns()[index].width
  }

  _getRowHeight = ({index}) => {
    const {lines} = this.state
    const line    = lines.get(index, Map())
    if (line.get('type') === 'filename') {
      return 42
    }

    const message   = line.get('msg', '')
    let lineCount

    if (Array.isArray(message)) {
      lineCount = message.reduce((count, token) => {
        if (_isString(token)) {
          return count + token.split('\n').length - 1
        } else {
          return count
        }
      }, 0)
    } else {
      lineCount = message.split('\n').length
    }

    return (lineCount > 1 ? lineCount + 2 : 1) * 16
  }

  _getColumnCount = () => this._getColumns().length

  _handleLineNumberClick = (lineNumber) => {
    const {location} = this.props
    const args       = queryString.parse(location.search)
    navigate(location.pathname, {...args, line: lineNumber})
  }

  _cellRenderer = ({
                     rowIndex: index,
                     columnIndex,
                     key,
                     style,
                   }) => {
    const {matches, findIndex, findQuery} = this.props
    const column                          = this._getColumns()[columnIndex]
    const line                            = this.state.lines.get(index)

    let content = null
    switch (column.name) {
      case 'msg': {
        content = line.get('msg')
        if (line.get('type') === 'filename') {
          content =
            <div className="file-name"><Link to={`/${content}?fs=${line.get('fs')}&line=${line.get('line')}`}>{content}</Link></div>
        }
        if (matches && matches.has(line.get('line'))) {
          const findIndexOffset = matches.keySeq().sort().reduce((result, matchIndex) => {
            const lineMatches = matches.get(matchIndex)
            if (matchIndex < line.get('line')
            ) {
              return result + lineMatches.length
            }
            return result
          }, 0)

          const tokens = content.split(findQuery)
          content      = _flatMap(tokens, (token, tokenIndex) => {
            const isCurrent = findIndexOffset + tokenIndex === findIndex
            return (tokenIndex < (tokens.length - 1) ? [token,
              <span
                className={cn('find-highlight', {current: isCurrent})}>{findQuery}</span>] : [token])
          })
        }
        break
      }
      case 'level': {
        content = line.get('level') ?
          <Tag key={line.get('level')} color={colorByLevel(line.get('level'))}>{line.get('level')}</Tag> : null
        break
      }
      case 'timestamp': {
        if (line.get('time')) {
          const timestamp = moment(line.get('time'))
          content         = <span className="time" key={index}>{timestamp.format('YY/MM/DD HH:mm:ss')}</span>
        }
        break
      }
      case 'linenumber': {
        content =
          <span className="linenumber" key={index} onClick={this._handleLineNumberClick.bind(this, line.get('line'))}>{line.get('line')}</span>
        break
      }
    }

    return (
      <div key={key} style={style} className={cn('line', line.get('type'), line.get('level', '').toLowerCase())}>
        {content}
      </div>
    )
  }

  _renderLines = () => {
    const {location = {}, scrollToLine} = this.props
    const {lines}                       = this.state

    const {line = 0} = queryString.parse(location.search)
    return (
      <AutoSizer>
        {({height, width}) => (
          <Grid
            ref={node => this.grid = node}
            scrollToRow={scrollToLine || Math.min(lines.size, Number(line))}
            width={width}
            height={height}
            rowCount={lines.size}
            rowHeight={this._getRowHeight}
            cellRenderer={this._cellRenderer}
            columnCount={this._getColumnCount()}
            columnWidth={this._getColumnWidth}
          />
        )}
      </AutoSizer>
    )
  }

  render() {
    return (
      <div className="lines-view-container">
        <div className="lines-view">
          {this._renderLines()}
        </div>
      </div>
    )
  }
}

export {LinesView}