import React, {Component} from 'react'
import {Input, Icon, Layout, Breadcrumb} from 'antd'
import FileTree from 'file-tree'
import {Route, Switch} from 'react-router'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {filesSelector, filterSelector, locationSelect} from 'selectors'
import {Link} from 'react-router-dom'
import FileView from 'file-view'
import queryString from 'query-string'
import {withLoader} from 'utils'
import {setFilter} from 'sockets/socket-actions'

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
    console.log('home.js@render: files', files)
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
    }
    return null
  }
}

class Home extends Component {
  render() {
    return (
      <Layout className="layout home">
        <Header>
          <div className="logo">Log Server</div>
          <Input
            placeholder="Search"
            className="search"
            prefix={<Icon type="search" style={{color: 'rgba(0,0,0,.25)'}}/>}
          />
        </Header>
        <Content style={{padding: '0 50px'}}>
          <Breadcrumbs/>
          <div style={{background: '#fff', padding: 24, minHeight: 280}}>
            <Switch>
              <Route path="/files/*" component={FileTree} exact={false}/>
              <Route path="/view" component={FileView} exact={true}/>
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
