import React, {Component} from 'react'
import _debounce from 'lodash/debounce'
import {Input, Icon, Layout, Breadcrumb} from 'antd'
import FileTree from 'file-tree'
import {Route, Switch} from 'react-router'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {filesSelector, hasPendingRequest, isIndexReady, locationSelect, searchSelector} from 'selectors'
import {Link, Redirect} from 'react-router-dom'
import queryString from 'query-string'
import {navigate, withLoader} from 'utils'
import {clearSearchResults, send, setFilter, setSearch} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'
import SearchView from 'file-view/search-view'
import Loader from 'loader/loader'

const {Header, Content, Footer} = Layout

@connect(createStructuredSelector({
  location:         locationSelect,
  files:            filesSelector,
  searchRequesting: hasPendingRequest(API_ACTIONS.SEARCH),
}), {
  setFilter,
})
class Breadcrumbs extends Component {
  state = {
    filter: '',
  }

  handleChange = (e) => {
    this.setState({
      filter: e.target.value,
    })
    this.updateFilter()
  }

  updateFilter = _debounce(() => {
    this.props.setFilter(this.state.filter)
  }, 300)

  render() {
    const {location, files, searchRequesting} = this.props
    const {search}                            = queryString.parse(location.search)

    if (search) {
      return (
        <Breadcrumb style={{margin: '16px 0'}} separator=">">
          <Breadcrumb.Item><Link to={'/'}>Home</Link></Breadcrumb.Item>
          <Breadcrumb.Item>Search Results {searchRequesting ? <Loader size={15}/> : null}</Breadcrumb.Item>
        </Breadcrumb>
      )
    } else {
      const path = ['Home'].concat(location.pathname.split('/').filter(Boolean))
      return (
        <Breadcrumb style={{margin: '16px 0'}} separator=">">
          {path.map((pathPart, i) => {
            return <Breadcrumb.Item key={pathPart}><Link to={`/${path.slice(1, i + 1).join('/')}`}>{pathPart}</Link></Breadcrumb.Item>
          })}
          {files.size ? <Breadcrumb.Item><input className="tree-search" placeholder="filter..."
                                                value={this.state.filter}
                                                onChange={this.handleChange}/>
          </Breadcrumb.Item> : null}
        </Breadcrumb>
      )
    }
  }
}

@connect(createStructuredSelector({
  search:       searchSelector,
  isIndexReady: isIndexReady,

}), {
  send,
  setSearch,
  clearSearchResults,
})
class Home extends Component {
  handleSearch = ({target: {value}}) => {
    this.props.setSearch(value)
    this.props.clearSearchResults()
    this.props.send(API_ACTIONS.SEARCH, {
      path:   [],
      regexp: value,
    })

    if (value) {
      const currentSearch = queryString.stringify({search: value})
      navigate(`/?${currentSearch}`)
    } else {
      navigate('/')
    }
  }

  componentDidMount() {
    const {location} = this.props
    const {search}   = queryString.parse(location.search)

    this.props.send(API_ACTIONS.GET_FILE_TREE, {
      base_path: [],
    })

    if (search) {
      this.handleSearch({target: {value: search}})
    }
  }

  componentWillReceiveProps(nextProps) {
    const {location} = nextProps
    if (!location.search) {
      this.props.setSearch('')
      this.props.clearSearchResults()
    }
  }

  _renderMainComponent = (props) => {
    const {location} = this.props
    const {search}   = queryString.parse(location.search)

    if (search) {
      return (
        <SearchView {...props}/>
      )
    }

    return (
      <FileTree {...props}/>
    )
  }

  render() {
    const {isIndexReady} = this.props

    const content = isIndexReady ? <Content style={{padding: '0 50px'}}>
      <Breadcrumbs/>
      <div className="main-content" style={{background: '#fff', padding: 24, minHeight: 280}}>
        <Switch>
          <Route path="/*" render={this._renderMainComponent} exact={false}/>
          <Redirect to={{
            pathname: '/',
          }}/>}/>
        </Switch>
      </div>
    </Content> : <Loader/>

    return (
      <Layout className="layout home">
        <Header>
          {this.props.requesting ? <div className="logo">REQUESTING</div> : <div className="logo">Log Server</div>}
          <Input
            placeholder="Search"
            className="search"
            value={this.props.search}
            onChange={this.handleSearch}
            prefix={<Icon type="search" style={{color: 'rgba(0,0,0,.25)'}}/>}
          />
        </Header>
        {content}
        <Footer style={{textAlign: 'center'}}>
        </Footer>
      </Layout>
    )
  }
}

export default withLoader(Home)
