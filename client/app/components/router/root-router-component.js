import React from 'react'
import {Provider} from 'react-redux'
import {browserHistory} from 'router/history'
import {ConnectedRouter} from 'react-router-redux'
import App from 'layouts/app'

export default class RootRouterComponent extends React.Component {
  render() {
    return (
      <Provider store={this.props.store}>
        <ConnectedRouter history={browserHistory}>
          <App/>
        </ConnectedRouter>
      </Provider>
    )
  }
}




