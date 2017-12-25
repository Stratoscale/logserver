import React, {Component} from 'react'
import {contentSelector, locationSelect} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import queryString from 'query-string'
import {send} from 'sockets/socket-actions'
import cn from 'classnames'
import {API_ACTIONS} from 'consts'
import {withLoader} from 'utils'


@connect(createStructuredSelector({
  location: locationSelect,
  content:  contentSelector,
}), {
  send,
})
class FileView extends Component {
  componentDidMount() {
    const {location, send} = this.props
    const search           = queryString.parse(location.search)

    send(API_ACTIONS.GET_CONTENT, {
      path: search.file.split('/').filter(Boolean)
    })
  }

  render() {
    const {content} = this.props
    console.log('file-view.js@render: content = ', content)
    const search = queryString.parse(this.props.location.search)
    return (
      <div>
        <div>{content.map((line, index) => <div key={index} className={cn(line.level)} >{line.msg}</div> )}</div>
      </div>
    )
  }
}

export default withLoader(FileView)