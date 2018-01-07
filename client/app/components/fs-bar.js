import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Tag} from 'antd'

const {CheckableTag} = Tag

class FSBar extends Component {
  static propTypes    = {
    items:    PropTypes.arrayOf(PropTypes.shape({
      name:   PropTypes.string,
      active: PropTypes.bool,
    })),
    onToggle: PropTypes.func.isRequired,
  }
  static defaultProps = {
    items: [],
  }

  render() {
    const {items} = this.props
    return (
      <div>
        {items.map((item, index) => <CheckableTag onChange={this.props.onToggle.bind(this, item.name)} key={index}
                                                  checked={item.active}>{item.name}</CheckableTag>)}
      </div>
    )
  }
}

export {FSBar}