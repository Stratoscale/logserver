import React, {Component} from 'react'
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

const calculateMaxLengths = (lines) => lines.reduce(({message, level}, line) => ({
  message: Math.max(message, line.get('msg').length),
  level:   Math.max(level, line.get('level').length),
}), {message: 0, level: 0})


class LinesView extends Component {
  state = {
    maxLengths: calculateMaxLengths(this.props.lines),
  }

  static propTypes = {
    showFilename:    PropTypes.bool,
    showTimestamp:   PropTypes.bool,
    showLinenumbers: PropTypes.bool,
    showLevels:      ImmutablePropTypes.set,
  }

  static defaultProps = {
    showFilename:    false,
    showTimestamp:   false,
    showLinenumbers: false,
    showLevels:      Set(),
  }

  componentWillReceiveProps({lines}) {
    if (!this.props.lines.equals(lines)) {
      this.setState({
        maxLengths: calculateMaxLengths(lines),
      })
    }
  }

  _renderWithFilename = () => {
    const {lines}      = this.props
    const groupedLines = lines.groupBy(line => line.get('file_name'))
    return (
      groupedLines.entrySeq().map(([filename, lines]) => {
        const firstLine = lines.first() || Map()
        return (
          <div className="file-results" key={filename}>
            <div className="file-name"><Link to={`/${filename}?fs=${firstLine.get('fs')}&line=${firstLine.get('line')}`}>{filename}</Link></div>
            {lines.take(5).map((line = Map(), index) => {
                return (
                  <div key={index} className={cn('line', line.get('level', '').toLowerCase())}>
                    {line.get('level') ? <Tag key={line.get('level')} color={colorByLevel(line.get('level'))}>{line.get('level')}</Tag> : null}
                    {line.get('msg')}
                  </div>
                )
              }
            )}
          </div>
        )
      })
    )
  }

  _getColumns = () => {
    const {maxLengths} = this.state

    const columns                          = [
      {
        width: maxLengths.level * 9,
        name:  'level',
      },
      {
        width: maxLengths.message * 8,
        name:  'msg',
      },
    ]
    const {showTimestamp, showLinenumbers} = this.props

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

  _getColumnWidth = ({index}) => this._getColumns()[index].width

  _getColumnCount = () => this._getColumns().length


  _renderWithoutFilename = () => {
    const {location = {}, showLevels} = this.props
    let {lines}                       = this.props
    if (showLevels.size > 0) {
      lines = lines.filter(line => !line.get('level') || showLevels.includes(line.get('level', '').toLowerCase()))
    }

    const {line = 0}      = queryString.parse(location.search)

    return (
      <AutoSizer>
        {({height, width}) => (
          <Grid
            scrollToRow={Number(line)}
            width={width}
            height={height}
            rowCount={lines.size}
            rowHeight={({index}) => {
              const line      = lines.getIn([index, 'msg'], '')
              const lineCount = line.split('\n').length
              return (lineCount > 1 ? lineCount + 2 : 1) * 16
            }}
            cellRenderer={({
                             rowIndex: index,
                             columnIndex,
                             key,
                             style,
                          }) => {
              const column = this._getColumns()[columnIndex]
              const line   = lines.get(index)
              let content
              switch (column.name) {
                case 'msg': {
                  content = line.get('msg')
                  break
                }
                case 'level': {
                  content = line.get('level') ?
                    <Tag key={line.get('level')} color={colorByLevel(line.get('level'))}>{line.get('level')}</Tag> : null
                  break
                }
                case 'timestamp': {
                  const timestamp = moment(line.get('time'))
                  content         = <span className="time" key={index}>{timestamp.format('YY/MM/DD HH:mm:ss')}</span>
                  break
                }
                case 'linenumber': {
                  content = <span className="linenumber" key={index}>{line.get('line')}</span>
                  break
                }
              }

              return (
                <div key={key} style={style} className={cn('line', line.get('level', '').toLowerCase())}>
                  {content}
                </div>
              )
            }}
            columnCount={this._getColumnCount()}
            columnWidth={this._getColumnWidth}
          />
        )}
      </AutoSizer>
    )
  }

  render() {
    const {showFilename} = this.props
    let content

    if (showFilename) {
      content = this._renderWithFilename()
    } else {
      content = this._renderWithoutFilename()
    }

    return (
      <div className="lines-view-container">
        <div className="lines-view">
          {content}
        </div>
      </div>
    )
  }
}

export {LinesView}