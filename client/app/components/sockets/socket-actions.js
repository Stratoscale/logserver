import {socket} from './index'
import {ACTIONS} from 'consts'

let messageId = 1

export function socketReady() {
  return {
    type: ACTIONS.SOCKET_READY,
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