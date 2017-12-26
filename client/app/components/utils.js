import {store} from 'store/store'
import {connect} from 'react-redux'
import {branch, compose, renderComponent} from 'recompose'
import {isSocketReady} from 'selectors'
import {createStructuredSelector} from 'reselect'
import {withRouter} from 'react-router-dom'
import Loader from 'loader/loader'
import {browserHistory} from 'router/history'

export const runSelector = selector => selector(store.getState())


export const withLoader = compose(
  withRouter,
  connect(
    createStructuredSelector({socket_ready: isSocketReady}),
  ),
  branch(
    ({socket_ready}) => !socket_ready,
    renderComponent(Loader)
  ),
)

export const queryParams = (args) => {
  const array = []
  for (let key of Object.keys(args)) {
    array.push(encodeURIComponent(key) + '=' + encodeURIComponent(args[key]))
  }
  return array.join('&')
}

export function navigate(route, query = {}, state = {}) {
  query = queryParams(query)
  browserHistory.push(route + (query ? `?${query}` : ''), state)
}
