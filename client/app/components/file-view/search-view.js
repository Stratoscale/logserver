import React, {Component} from 'react'
import {filesSelector, locationSelect, searchResultsSelector} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {send} from 'sockets/socket-actions'
import {LinesView} from 'file-view/lines-view'

@connect(createStructuredSelector({
  location:    locationSelect,
  results:     searchResultsSelector,
  files:       filesSelector,
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