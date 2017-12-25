import React, {Component} from 'react'
import {Layout, Menu, Breadcrumb} from 'antd'
import SocketContainer from 'sockets'
import FileTree from 'file-tree'
import {Route, Switch} from 'react-router'
import {connect} from 'react-redux'
import {createStructuredSelector} from 'reselect'
import {locationSelect} from 'selectors'
import {Link} from 'react-router-dom'

const {Header, Content, Footer} = Layout

@connect(createStructuredSelector({
  location: locationSelect,
}))
class Breadcrumbs extends Component {
  render() {
    const path = this.props.location.pathname.split('/').filter(Boolean)
    return (
      <Breadcrumb style={{margin: '16px 0'}}>
        {path.map((pathPart, i) => {
          return <Breadcrumb.Item key={pathPart}><Link to={`/files/${path.slice(1, i + 1).join('/')}`}>{pathPart}</Link></Breadcrumb.Item>
        })}
      </Breadcrumb>
    )
  }
}

class Home extends Component {
  render() {
    return (
      <Layout className="layout home">
        <SocketContainer/>
        <Header>
          <div className="logo">Log Server/Streamer/Viewer</div>
        </Header>
        <Content style={{padding: '0 50px'}}>
          <Breadcrumbs/>
          <div style={{background: '#fff', padding: 24, minHeight: 280}}>
            <Switch>
              <Route path="/files/*" component={FileTree} exact={false}/>
              <Route path="/files" component={FileTree} exact={false}/>
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


export default Home
