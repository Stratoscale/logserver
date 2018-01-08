import React, {Component} from 'react'
import cn from 'classnames'
import {Map} from 'immutable'
import {Tag} from 'antd'
import PropTypes from 'prop-types'
import {List, WindowScroller, AutoSizer} from 'react-virtualized'


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
            <div className="file-name">{filename}</div>
            {lines.map((line = Map(), index) => {
                return (
                  <div key={index} className={cn('line', line.get('level', '').toLowerCase())}>
                    {line.get('level') ? <Tag color={colorByLevel(line.get('level'))}>{line.get('level')}</Tag> : null}
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
    const {lines} = this.props
    return (
      <AutoSizer>
        {({height, width}) => (
          <List
            width={width}
            height={height}
            rowCount={lines.size}
            rowHeight={({index}) => {
              const lineSize = lines.getIn([index, 'msg'], '').length
              return Math.ceil(lineSize * 13 / width) * 16
            }}
            rowRenderer={({
                            index,
                            key,
                            style
                          }) => {
              const line    = lines.get(index)
              const content = [line.get('level') ?
                <Tag color={colorByLevel(line.get('level'))}>{line.get('level')}</Tag> : null, line.get('msg')]
              return (
                <div key={key} style={style} className={cn('line', line.get('level', '').toLowerCase())}>
                  {content}
                </div>
              )
            }}
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