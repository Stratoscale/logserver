import React from 'react'
import {Provider} from 'react-redux'
import {browserHistory} from 'router/history'
import {Route, Switch} from 'react-router'
import {ConnectedRouter} from 'react-router-redux'
import FolderTree from 'folder-tree';

export default class RootRouterComponent extends React.Component {
  render() {
    return (
      <Provider store={this.props.store}>
        <ConnectedRouter history={browserHistory}>
          <Switch>
            <Route path="/" component={FolderTree}/>
          </Switch>
        </ConnectedRouter>
      </Provider>
    );
  }
}




