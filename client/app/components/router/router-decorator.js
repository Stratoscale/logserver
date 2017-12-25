import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import {locationSelector} from 'router/router-selectors'

export const navigation = () => (ComponentToInject) => {
  return connect(locationSelector, {navigate: (pathname, search) => push({pathname, search})}, null, {withRef: true})(ComponentToInject);
}