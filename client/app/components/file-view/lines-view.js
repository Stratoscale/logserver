import React, {Component} from 'react'
import cn from 'classnames'
import {Map} from 'immutable'
import {Tag} from 'antd'

const {CheckableTag} = Tag

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

class MyTag extends Component {
  constructor(props) {
    super(props)
    this.state = {
      checked: true,
    }
  }

  handleChange = (checked) => {
    this.setState({checked})
  }

  render() {
    return <CheckableTag {...this.props} checked={this.state.checked} onChange={this.handleChange}/>
  }
}

class LinesView extends Component {
  render() {
    const {lines} = this.props
    return (
      <div className="lines-view-container">
        <div className="controls">
          <MyTag color="blue">INFO</MyTag>
          <MyTag color="gold">WARNING</MyTag>
          <MyTag color="red">ERROR</MyTag>
        </div>
        <div className="lines-view">
          {lines.map((line = Map(), index) => <div key={index} className={cn('line', line.get('level', '').toLowerCase())}>
            {line.get('level') ? <Tag color={colorByLevel(line.get('level'))}>{line.get('level')}</Tag> : null} {line.get('msg')}</div>)}
        </div>
      </div>
    )
  }
}

export default LinesView