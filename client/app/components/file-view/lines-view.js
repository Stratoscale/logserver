import React, {Component} from 'react'
import cn from 'classnames'
import {Map} from 'immutable'
import {Tag} from 'antd'
import PropTypes from 'prop-types'
import {Grid, AutoSizer} from 'react-virtualized'
import {Link} from 'react-router-dom'
import queryString from 'query-string'


const colorByLevel = (level = '') => {
  switch (level.toLowerCase()) {
    case 'info':
      return 'blue'
    case 'error':
      return 'red'
    case 'warning':
      return 'gold'
  }
}

class LinesView extends Component {
  static propTypes = {
    showFilename: PropTypes.bool,
  }

  static defaultProps = {
    showFilename: false,
  }

  _renderWithFilename = () => {
    const {lines}      = this.props
    const groupedLines = lines.groupBy(line => line.get('file_name'))
    return (
      groupedLines.entrySeq().map(([filename, lines]) => {
        return (
          <div className="file-results" key={filename}>
            <div className="file-name"><Link to={`/${filename}?line=${lines.first().get('line')}`}>{filename}</Link></div>
            {lines.map((line = Map(), index) => {
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

  _renderWithoutFilename = () => {
    const {lines}         = this.props
    const {location = {}} = this.props
    const {line = 0}      = queryString.parse(location.search)
    const maxLineLength   = lines.reduce((result, line) => {
      return Math.max(result, line.get('msg').length)
    }, 0)

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
                            key,
                            style,
                          }) => {
              const line    = lines.get(index)
              const content = [line.get('level') ?
                <Tag key={line.get('level')} color={colorByLevel(line.get('level'))}>{line.get('level')}</Tag> : null, line.get('msg')]
              return (
                <div key={key} style={style} className={cn('line', line.get('level', '').toLowerCase())}>
                  {content}
                </div>
              )
            }}
            columnCount={1}
            columnWidth={maxLineLength * 8}
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