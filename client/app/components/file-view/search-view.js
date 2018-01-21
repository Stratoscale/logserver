import React, {Component} from 'react'
import {filesSelector, hasPendingRequest, locationSelect, searchResultsSelector} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {send} from 'sockets/socket-actions'
import {LinesView} from 'file-view/lines-view'
import {API_ACTIONS} from 'consts'

@connect(createStructuredSelector({
  location:   locationSelect,
  results:    searchResultsSelector,
  files:      filesSelector,
  requesting: hasPendingRequest(API_ACTIONS.SEARCH),
}), {
  send,
})
class SearchView extends Component {
  render() {
    const {results, ...props} = this.props

    if (!results.size) {
      return (
        <div>No results found</div>
      )
    }

    return (
      <div className="search-view">
        <LinesView {...props} lines={results} showFilename/>
      </div>
    )
  }
}

export default SearchView