import {createReducer} from 'redux-immutablejs';
import {fromJS, Map, List} from 'immutable';
import {ACTIONS} from 'consts';

const INITIAL_STATE  = fromJS({});
const INITIAL_ENTITY = fromJS([]);

export function addEntity(state, action) {
  const entityKey  = action.meta.key;
  const entityType = action.meta.type;

  const entity = fromJS(action.payload || INITIAL_ENTITY);
  if (entity.equals(state.getIn([entityType, entityKey]))) {
    return state
  } else if (Array.isArray(entityKey)) {
    return state.setIn([entityType, ...entityKey], entity);
  } else {
    return state.setIn([entityType, entityKey], entity);
  }
}

export function setEntities(state = INITIAL_STATE, action) {
  if (action.payload) {
    if (action.meta.merge) {
      if (action.meta.keepAPIOrder) {
        const type = state.get(action.meta.type, Map());

        const items = type.get('items', Map()).mergeDeep(fromJS(action.payload.items));
        const ids   = type.get('ids', List()).concat(List(action.payload.ids));

        return state.set(action.meta.type, type.withMutations(type => {
          type.set('items', items);
          type.set('ids', ids);
        }));

      } else {
        return state.set(action.meta.type, state.get(action.meta.type, Map()).mergeDeep(fromJS(action.payload)));
      }
    } else if (action.meta.receivedOnly) {
      const current = state.get(action.meta.type, Map());
      const payload = fromJS(action.payload);
      return state.set(action.meta.type, current.withMutations(typeState => {
        for (let [key, value] of payload.entries()) {
          typeState.set(key, value);
        }
      }));
    } else {
      return state.set(action.meta.type, fromJS(action.payload));
    }
  }
  return state;
}

export function deleteEntity(state, action) {
  const {type, id} = action.payload;
  return state.deleteIn([type, id]);
}

export function deleteEntities(state, action) {
  const {type} = action.payload;
  return state.delete(type);
}

export function clearEntities(state, action) {
  const keepInStore = (entities) =>
    INITIAL_STATE.withMutations(
      (stateWithEntities) =>
        entities.map((entityKey) =>
          stateWithEntities.set(entityKey, state.get(entityKey)))
    );

  return action.meta.keep ? keepInStore(action.meta.keep) : INITIAL_STATE
}

export function initEntities(state, {meta: {type}}) {
  if (!state.get(type)) {
    return state.set(type, Map())
  }
  return state;
}

export const entities = createReducer(INITIAL_STATE, {
  [ACTIONS.ADD_ENTITY]:          addEntity,
  [ACTIONS.SET_ENTITIES]:        setEntities,
  [ACTIONS.DELETE_ENTITY]:       deleteEntity,
  [ACTIONS.DELETE_ENTITIES]:     deleteEntities,
  [ACTIONS.CLEAR_ENTITIES]:      clearEntities,
  [ACTIONS.FAIL_REQUEST_ENTITY]: initEntities, // Make sure the type is initialized even when requests fail, so that wereFetched selectors work correctly
});


export default entities;