import keyMirror from 'keymirror'

export const ACTIONS = keyMirror({
  SOCKET_READY: null,
  SET_FILES:    null,
})

export const API_ACTIONS = {
  GET_FILE_TREE: 'get-file-tree',
}
