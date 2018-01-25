import React, {Component} from 'react'
import _debounce from 'lodash/debounce'
import {filesSelector, findSelector, hasPendingRequest, indexSelector, locationSelect, matchesSelector} from 'selectors'
import {API_ACTIONS} from 'consts'
import {createStructuredSelector} from 'reselect'
import {clearFind, findNext, findPrev, setFilter, setFind} from 'sockets/socket-actions'
import {connect} from 'react-redux'
import queryString from 'query-string'
import {Breadcrumb, Button, Input} from 'antd'
import Loader from 'loader/loader'
import {Link} from 'react-router-dom'

@connect(createStructuredSelector({
  index:            indexSelector,
  location:         locationSelect,
  find:             findSelector,
  matches:          matchesSelector,
  files:            filesSelector,
  searchRequesting: hasPendingRequest(API_ACTIONS.SEARCH),
}), {
  setFilter,
  setFind,
  findNext,
  findPrev,
  clearFind,
})
class Breadcrumbs extends Component {
  state = {
    filter: '',
    find:   '',
  }

  handleChange = (e) => {
    this.setState({
      filter: e.target.value,
    })
    this.updateFilter()
  }

  handleFindChange = (e) => {
    const find = e.target.value
    if (find) {
      this.setState(() => ({
        find,
      }), () => {
        this.updateFind({
          query: this.state.find,
        })
      })
    } else {
      this._clearFind()
    }
  }

  _nextFindResult = () => {
    this.props.findNext()
  }

  _prevFindResult = () => {
    this.props.findPrev()
  }

  _clearFind = () => {
    this.setState(() => ({
      find: '',
    }), () => this.props.clearFind())

  }

  updateFilter = _debounce(() => {
    this.props.setFilter(this.state.filter)
  }, 300)

  updateFind = _debounce(() => {
    this.props.setFind(this.state.find)
  }, 300)

  _getFilterComponent = () => {
    const {index, location: {pathname}, find, matches} = this.props

    const isDir = pathname.length === 1 || index.getIn([pathname.substr(1), 'is_dir'])

    if (isDir) {
      return (
        <Breadcrumb.Item>
          <Input className="tree-search"
                 ref={this._inputRef}
                 placeholder="filter..."
                 value={this.state.filter}
                 onChange={this.handleChange}
          />
        </Breadcrumb.Item>
      )
    } else {
      const totalMatches = matches.valueSeq().reduce((result, lineMatches) => result + lineMatches.length, 0)
      return (
        <Breadcrumb.Item>
          <Input className="tree-search"
                 ref={this._inputRef}
                 placeholder="find in file..."
                 addonAfter={<Button.Group>
                   <Button icon="up" onClick={this._prevFindResult}/>
                   <Button icon="down" onClick={this._nextFindResult}/>
                   <Button icon="close" onClick={this._clearFind}/>
                 </Button.Group>}
                 value={this.state.find}
                 onChange={this.handleFindChange}
          />
          {find.get('query') ? matches.size ? <span
              className="find-index">{Math.min(find.get('index', 0) + 1, totalMatches)}/{totalMatches}</span> :
            <span className="find-index">No Match</span> : null}

        </Breadcrumb.Item>
      )

    }
  }

  _inputRef = node => this.input = node

  _keyDownHandler = (e) => {
    if (e.key === 'f' && e.ctrlKey === true) {
      e.preventDefault()
      if (this.input) {
        this.input.focus()
      }
    }
  }

  componentDidMount() {
    window.addEventListener('keydown', this._keyDownHandler)
  }

  componentWillUnmount() {
    window.removeEventListener('keydown', this._keyDownHandler)
  }

  render() {
    const {location: {pathname}, searchRequesting} = this.props

    const {search} = queryString.parse(location.search)

    if (search) {
      return (
        <Breadcrumb style={{margin: '10px 0'}} separator=">">
          <Breadcrumb.Item><Link to={'/'}>Home</Link></Breadcrumb.Item>
          <Breadcrumb.Item>Search Results {searchRequesting ? <Loader size={15}/> : null}</Breadcrumb.Item>
        </Breadcrumb>
      )
    } else {
      const path = ['Home'].concat(pathname.split('/').filter(Boolean))
      return (
        <Breadcrumb style={{margin: '10px 0'}} separator=">">
          {path.map((pathPart, i) => {
            return <Breadcrumb.Item key={pathPart}><Link to={`/${path.slice(1, i + 1).join('/')}`}>{pathPart}</Link></Breadcrumb.Item>
          })}
          {this._getFilterComponent()}
        </Breadcrumb>
      )
    }
  }
}

export default Breadcrumbs