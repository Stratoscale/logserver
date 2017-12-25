import {socket} from './index'
import {ACTIONS} from 'consts'

let messageId = 0

export function socketReady() {
  return {
    type: ACTIONS.SOCKET_READY,
  }
}

export function send(action) {
  return (dispatch, getState) => {
    socket.send(JSON.stringify({
        meta:      {
          action,
          id: messageId++,
        },
        base_path: [],
      })
    )
  }
}