import React, {Component} from 'react'
import _noop from 'lodash/noop'
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
  message: Math.max(message, line.get('msg').length),
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
  } else if (showLevels.size > 0) {
    return lines.filter(line => !line.get('level') || showLevels.includes(line.get('level', '').toLowerCase()))
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
    showLevels:      ImmutablePropTypes.set,
    scrollToLine:    PropTypes.number,
    onScroll:        PropTypes.func,
  }

  static defaultProps = {
    showFilename:    false,
    showTimestamp:   false,
    showLinenumbers: false,
    showLevels:      Set(),
    scrollToLine:    0,
    onScroll:        _noop,
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
    const lineCount = message.split('\n').length
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
    const column = this._getColumns()[columnIndex]
    const line   = this.state.lines.get(index)

    let content = null
    switch (column.name) {
      case 'msg': {
        content = line.get('msg')
        if (line.get('type') === 'filename') {
          content =
            <div className="file-name"><Link to={`/${content}?fs=${line.get('fs')}&line=${line.get('line')}`}>{content}</Link></div>
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
    const {location = {}, onScroll, scrollToLine} = this.props
    const {lines}                                 = this.state

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
            onScroll={({scrollTop}) => {
              onScroll(Math.round(scrollTop / 8))
            }}
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