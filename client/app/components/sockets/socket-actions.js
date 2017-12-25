import {socket} from './index'
import {ACTIONS} from 'consts'

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

export function setContent(payload) {
  return {
    type: ACTIONS.SET_CONTENT,
    payload,
  }
}

export function send(action, data) {
  return (dispatch, getState) => {
    socket.send(JSON.stringify({
        meta: {
          action,
          id: messageId++,
        },
        ...data,
      })
    )
  }
}