import createBrowserHistory from 'history/createBrowserHistory'

export const browserHistory = createBrowserHistory({
  basename: window.__INIT__.basePath,
})


