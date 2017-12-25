import React from 'react'
import {Provider} from 'react-redux'
import {browserHistory} from 'router/history'
import {Route, Switch} from 'react-router'
import {ConnectedRouter} from 'react-router-redux'
import Home from 'layouts/home'

export default class RootRouterComponent extends React.Component {
  render() {
    return (
      <Provider store={this.props.store}>
        <ConnectedRouter history={browserHistory}>
          <Switch>
            <Route path="/" component={Home}/>
          </Switch>
        </ConnectedRouter>
      </Provider>
    );
  }
}




