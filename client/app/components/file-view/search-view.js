import React, {Component} from 'react'
import {filesSelector, fileSystemsSelector, locationSelect, searchResultsSelector} from 'selectors'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {send} from 'sockets/socket-actions'
import LinesView from 'file-view/lines-view'

@connect(createStructuredSelector({
  location:    locationSelect,
  results:     searchResultsSelector,
  files:       filesSelector,
  fileSystems: fileSystemsSelector,
}), {
  send,
})
class SearchView extends Component {
  render() {
    const {results, fileSystems} = this.props


    if (!results.size) {
      return (
        <div>No results found</div>
      )
    }

    return (
      <div>
        {fileSystems}
        <LinesView lines={results}/>
      </div>
    )
  }
}

export default SearchView