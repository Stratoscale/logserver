import React, {Component} from 'react'
import cn from 'classnames'
import {Tag} from 'antd'

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
  render() {
    const {lines} = this.props
    return (
      <div className="lines-view">
        <div className="controls">Controls</div>
        {lines.map((line, index) => <div key={index} className={cn('line', line.get('level', '').toLowerCase())}><Tag
          color={colorByLevel(line.get('level'))}>{line.get('level')}</Tag> {line.get('msg')}</div>)}
      </div>
    )
  }
}

export default LinesView