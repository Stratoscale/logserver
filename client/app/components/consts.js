import keyMirror from 'keymirror'

export const ACTIONS = keyMirror({
  CLEAR_CONTENT:      null,
  CLEAR_FIND:         null,
  INDEX_READY:        null,
  RECEIVE_REQUEST:    null,
  FIND_NEXT:          null,
  FIND_PREV:          null,
  SEND_REQUEST:       null,
  SET_CONTENT:        null,
  SET_CONTENT_ID:     null,
  SET_CURRENT_PATH:   null,
  SET_FILES:          null,
  SET_FILTER:         null,
  SET_FIND:           null,
  SET_SEARCH:         null,
  SET_SEARCH_ID:      null,
  SET_SEARCH_RESULTS: null,
  SET_LEVELS:         null,
  SOCKET_READY:       null,
})

export const API_ACTIONS = {
  GET_FILE_TREE: 'get-file-tree',
  GET_CONTENT:   'get-content',
  SEARCH:        'search',
}

export const colorByLevel = (level = '') => {
  switch (level.toLowerCase()) {
    case 'debug':
      return 'cyan'
    case 'info':
      return 'blue'
    case 'error':
      return 'red'
    case 'warning':
      return 'gold'
    case 'success':
      return 'green'
    case 'progress':
      return 'purple'
  }
}
