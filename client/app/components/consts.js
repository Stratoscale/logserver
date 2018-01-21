import keyMirror from 'keymirror'

export const ACTIONS = keyMirror({
  SOCKET_READY:       null,
  SET_FILES:          null,
  SET_CURRENT_PATH:   null,
  SET_CONTENT:        null,
  CLEAR_CONTENT:      null,
  SET_CONTENT_ID:     null,
  SET_FILTER:         null,
  SET_SEARCH:         null,
  SET_SEARCH_ID:      null,
  SET_SEARCH_RESULTS: null,
  SEND_REQUEST:       null,
  RECEIVE_REQUEST:    null,
})

export const API_ACTIONS = {
  GET_FILE_TREE: 'get-file-tree',
  GET_CONTENT:   'get-content',
  SEARCH:        'search',
}
