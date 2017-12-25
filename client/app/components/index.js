import React from 'react'
import ReactDOM from 'react-dom';
import RootRouterComponent from 'router/root-router-component';
import {store} from 'store/store'
import {AppContainer} from 'react-hot-loader'

console.log('root mount');
const rootEl = document.querySelector('#root');

const App = ({isHotReload, RootComponent}) => {
  if (isHotReload) {
    require('react-hot-loader/patch')
    return (
      <AppContainer>
        <RootComponent store={store}/>
      </AppContainer>
    )
  }
  else
    return (<RootRouterComponent store={store}/>)
}

if (__DEV__) {
  ReactDOM.render(App({isHotReload: true, RootComponent: RootRouterComponent}), rootEl);

  // Hot Module Replacement API
  if (module.hot) {
    module.hot.accept('router/root-router-component', () => {
      const NextApp = require('router/root-router-component').default;
      ReactDOM.render(App({isHotReload: true, RootComponent: NextApp}),
        rootEl
      );
    });
  }
}

else {
  ReactDOM.render(App({isHotReload: false}), rootEl);
}
