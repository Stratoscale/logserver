import {store} from 'store/store'
import {connect} from 'react-redux'
import {branch, compose, renderComponent} from 'recompose'
import {isSocketReady} from 'selectors'
import {createStructuredSelector} from 'reselect'
import Loader from 'loader/loader'

export const runSelector = selector => selector(store.getState())


export const withLoader = compose(
  connect(
    createStructuredSelector({socket_ready: isSocketReady}),
  ),
  branch(
    ({socket_ready}) => !socket_ready,
    renderComponent(Loader)
  )
)