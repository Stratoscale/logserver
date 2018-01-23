import {Component} from 'react'
import {connect} from 'react-redux'
import {addSearchResults, indexReady, receiveRequest, setContent, setFiles, socketReady} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'

let socket = null

@connect(null, {socketReady, indexReady, setFiles, setContent, addSearchResults, receiveRequest})
export default class SocketContainer extends Component {
  constructor(props) {
    super(props)
  }

  componentWillMount() {
    if (!socket) {
      socket = new WebSocket(`ws://${window.location.host}${window.__INIT__.basePath}/_ws`)

      // Connection opened
      socket.addEventListener('open', (event) => {
        this.props.socketReady()
        console.log('socket opened')
      })

      // Listen for messages
      socket.addEventListener('message', (event) => {
        const {meta, finished, ...payload} = JSON.parse(event.data)
        if (finished) {
          this.props.receiveRequest(meta.action)
          if (meta.action === API_ACTIONS.GET_FILE_TREE) {
            this.props.indexReady()
          }
        }

        switch (meta.action) {
          case API_ACTIONS.GET_FILE_TREE: {
            this.props.setFiles(payload.tree)
            break
          }
          case API_ACTIONS.GET_CONTENT: {
            this.props.setContent(payload.lines || [], meta.id)
            break
          }
          case API_ACTIONS.SEARCH: {
            if (payload.lines) {
              this.props.addSearchResults(payload.lines, meta.id)
            }
            break
          }
          default: {
            console.warn('Unknown action returned from API', meta, payload)
          }
        }

      })
    }
  }

  render = () => {
    return null
  }
}

export {
  socket,
}