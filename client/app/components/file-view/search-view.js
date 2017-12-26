import React, {Component} from 'react'
import {filesSelector, locationSelect, searchResultsSelector} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {send} from 'sockets/socket-actions'
import LinesView from 'file-view/lines-view'

@connect(createStructuredSelector({
  location: locationSelect,
  results:  searchResultsSelector,
  files:    filesSelector,
}), {
  send,
})
class SearchView extends Component {
  render() {
    const {results} = this.props
    if (!results.size) {
      return (
        <div>No results found</div>
      )
    }

    return (
      <div>
        <LinesView lines={results.toJS()}/>
      </div>
    )
  }
}

export default SearchView