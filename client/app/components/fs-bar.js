import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Tag} from 'antd'
import classNames from 'classnames'

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
    const {items, className} = this.props
    return (
      <div className={classNames('fs-bar', className)}>
        {items.map((item, index) => <CheckableTag onChange={this.props.onToggle.bind(this, item.name)}
                                                  key={index}
                                                  checked={item.active}
                                                  color={item.color}
                                                  className={item.active ? `ant-tag-${item.color}` : ''}
        >{item.content || item.name}</CheckableTag>)}
      </div>
    )
  }
}

export {FSBar}