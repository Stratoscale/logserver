import React, {Component} from 'react'
import {Input, Icon, Layout, Breadcrumb} from 'antd'
import FileTree from 'file-tree'
import {Route, Switch} from 'react-router'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {filesSelector, filterSelector, locationSelect, searchSelector} from 'selectors'
import {Link, Redirect} from 'react-router-dom'
import FileView from 'file-view'
import queryString from 'query-string'
import {navigate, withLoader} from 'utils'
import {clearSearchResults, send, setFilter, setSearch} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'
import SearchView from 'file-view/search-view'

const {Header, Content, Footer} = Layout

@connect(createStructuredSelector({
  location: locationSelect,
  files:    filesSelector,
  filter:   filterSelector,
}), {
  setFilter,
})
class Breadcrumbs extends Component {
  handleChange = (e) => {
    this.props.setFilter(e.target.value)

  }

  render() {
    const {location, files, filter} = this.props
    if (location.pathname.startsWith('/view')) {
      const filename = queryString.parse(location.search).file
      return (
        <Breadcrumb style={{margin: '16px 0'}} separator=">">
          <Breadcrumb.Item><Link to={'/files/'}>Home</Link></Breadcrumb.Item>
          <Breadcrumb.Item>{filename}</Breadcrumb.Item>
        </Breadcrumb>
      )
    } else if (location.pathname.startsWith('/files')) {
      const path = location.pathname.split('/').filter(Boolean).map(item => item === 'files' ? 'Home' : item)
      return (
        <Breadcrumb style={{margin: '16px 0'}} separator=">">
          {path.map((pathPart, i) => {
            return <Breadcrumb.Item key={pathPart}><Link to={`/files/${path.slice(1, i + 1).join('/')}`}>{pathPart}</Link></Breadcrumb.Item>
          })}
          {files.size ? <Breadcrumb.Item><input className="tree-search" placeholder="filter..." value={filter}
                                                onChange={this.handleChange}/>
          </Breadcrumb.Item> : null}
        </Breadcrumb>
      )

    } else if (location.pathname.startsWith('/search')) {
      return (
        <Breadcrumb style={{margin: '16px 0'}} separator=">">
          <Breadcrumb.Item><Link to={'/files/'}>Home</Link></Breadcrumb.Item>
          <Breadcrumb.Item>Search Results</Breadcrumb.Item>
        </Breadcrumb>
      )
    }
    return null
  }
}

@connect(createStructuredSelector({
  search: searchSelector,
}), {
  send,
  setSearch,
  clearSearchResults,
})
class Home extends Component {
  handleSearch = (e) => {
    this.props.setSearch(e.target.value)
    this.props.clearSearchResults()
    this.props.send(API_ACTIONS.SEARCH, {
      path:   [],
      regexp: e.target.value,
    })
    if (e.target.value) {
      navigate('/search')
    } else {
      navigate('/')
    }
  }

  render() {
    return (
      <Layout className="layout home">
        <Header>
          <div className="logo">Log Server</div>
          <Input
            placeholder="Search"
            className="search"
            value={this.props.search}
            onChange={this.handleSearch}
            prefix={<Icon type="search" style={{color: 'rgba(0,0,0,.25)'}}/>}
          />
        </Header>
        <Content style={{padding: '0 50px'}}>
          <Breadcrumbs/>
          <div style={{background: '#fff', padding: 24, minHeight: 280}}>
            <Switch>
              <Route path="/files/*" component={FileTree} exact={false}/>
              <Route path="/view" component={FileView} exact={true}/>
              <Route path="/search" component={SearchView} exact={true}/>
              <Redirect to={{
                pathname: '/files/',
              }}/>}/>
            </Switch>
          </div>
        </Content>
        <Footer style={{textAlign: 'center'}}>
          StratoHackathon 2017
        </Footer>
      </Layout>
    )
  }
}


export default withLoader(Home)
