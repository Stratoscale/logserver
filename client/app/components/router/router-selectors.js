import {createSelector, createStructuredSelector} from 'reselect';
import {Map} from 'immutable'

export const routingSelector           = createSelector((state) => state, (state = Map()) => {
  return state.get('routing', {}).locationBeforeTransitions || {};
});
export const routingStructuredSelector = createStructuredSelector({routing: routingSelector});

export const locationSelect = (state = Map()) => state.get('router').location || {}

export const locationSelector = createStructuredSelector({location: locationSelect})

export default {
  routingSelector,
  routingStructuredSelector,
};
