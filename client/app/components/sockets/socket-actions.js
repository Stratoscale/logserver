import {socket} from './index'
import {List} from 'immutable'
import {ACTIONS, API_ACTIONS} from 'consts'

let messageId = 1

export function socketReady() {
  return {
    type: ACTIONS.SOCKET_READY,
  }
}

export function setFiles(files) {
  return {
    type:    ACTIONS.SET_FILES,
    payload: files,
  }
}

export function setContent(payload, id) {
  return {
    type: ACTIONS.SET_CONTENT,
    meta: {
      id,
    },
    payload,
  }
}

export function setFilter(payload) {
  return {
    type: ACTIONS.SET_FILTER,
    payload,
  }
}

export function setSearch(payload) {
  return {
    type: ACTIONS.SET_SEARCH,
    payload,
  }
}

export function addSearchResults(payload, id) {
  return {
    type: ACTIONS.SET_SEARCH_RESULTS,
    meta: {
      add: true,
      id,
    },
    payload,
  }
}

export function clearSearchResults() {
  return {
    type:    ACTIONS.SET_SEARCH_RESULTS,
    meta:    {},
    payload: List(),
  }
}

export function setSearchId(id) {
  return {
    type:    ACTIONS.SET_SEARCH_ID,
    payload: id,
  }
}

export function setContentId(id) {
  return {
    type:    ACTIONS.SET_CONTENT_ID,
    payload: id,
  }
}

export function send(action, data) {
  const thunk = (dispatch, getState) => {
    const id = messageId++
    if (action === API_ACTIONS.SEARCH) {
      dispatch(setSearchId(id))
    }
    if (action === API_ACTIONS.GET_CONTENT) {
      dispatch(setContentId(id))
    }
    socket.send(JSON.stringify({
        meta: {
          action,
          id,
        },
        ...data,
      })
    )
  }
  thunk.meta  = {
    debounce: {
      time: 300,
      key:  'send-' + action,
    },
  }

  return thunk

}