import React, {Component} from 'react'
import {Layout, Menu, Breadcrumb} from 'antd'
import SocketContainer from 'sockets'
import FileTree from 'file-tree'
import {Route, Switch} from 'react-router'

const {Header, Content, Footer} = Layout

class Home extends Component {
  render() {
    return (
      <Layout className="layout home">
        <SocketContainer/>
        <Header>
          <div className="logo">Log Server/Streamer/Viewer</div>
        </Header>
        <Content style={{padding: '0 50px'}}>
          <Breadcrumb style={{margin: '16px 0'}}>
            <Breadcrumb.Item>Files</Breadcrumb.Item>
          </Breadcrumb>
          <div style={{background: '#fff', padding: 24, minHeight: 280}}>
            <Switch>
              <Route path="/" component={FileTree}/>
            </Switch>
          </div>
        </Content>
        <Footer style={{textAlign: 'center'}}>
          Ant Design ©2016 Created by Ant UED
        </Footer>
      </Layout>
    )
  }
}


export default Home
