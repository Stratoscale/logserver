/*global describe, it, beforeEach, assert*/
import {tokenValidation, __RewireAPI__ as createStoreRewire} from './create-store';
import {expect} from 'chai';
import {ACTIONS} from 'entities/entities-constants';
import {fromJS} from 'immutable';
import configureMockStore from 'redux-mock-store'
import thunk from 'redux-thunk'

const middlewares = [thunk];
const mockStore   = configureMockStore(middlewares);
const store = mockStore(fromJS({}));
createStoreRewire.__Rewire__('validateToken', () => ({type: 'validateTokenDispatched'}));

const dispatchWithStoreOf = (action) => {
  // eslint-disable-next-line no-unused-vars
  let dispatched = null;
  const dispatch = tokenValidation(store)(actionAttempt => dispatched = actionAttempt);
  dispatch(action);
  return store.getActions();
};

describe('Test tokenValidation middleware', () => {
  it('should dispatch if action is a fail request', () => {
    store.clearActions();
    const action = {
      type: ACTIONS.FAIL_REQUEST_ENTITY,
      meta: {},
    };

    expect(
      dispatchWithStoreOf(action)
    ).to.deep.equal([{type: 'validateTokenDispatched'}]);
  });

  it('should not dispatch if action is anything but a failed request', () => {
    store.clearActions();
    const action = {
      type: ACTIONS.SUCCESS_REQUEST_ENTITY,
      meta: {},
    };

    expect(
      dispatchWithStoreOf(action)
    ).to.deep.equal([]);
  })
});