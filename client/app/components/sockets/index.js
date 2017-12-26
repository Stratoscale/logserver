import {Component} from 'react'
import {connect} from 'react-redux'
import {setContent, setFiles, socketReady} from 'sockets/socket-actions'
import {API_ACTIONS} from 'consts'

let socket = null

@connect(null, {socketReady, setFiles, setContent})
export default class SocketContainer extends Component {
  constructor(props) {
    super(props)
  }

  componentWillMount() {
    if (!socket) {

      socket = new WebSocket(`ws://${window.location.host}/ws`)

      // Connection opened
      socket.addEventListener('open', (event) => {
        this.props.socketReady()
        console.log('socket opened')
      })

      // Listen for messages
      socket.addEventListener('message', (event) => {
        const {meta, ...payload} = JSON.parse(event.data)
        switch (meta.action) {
          case API_ACTIONS.GET_FILE_TREE: {
            this.props.setFiles(payload.tree)
            break;
          }
          case API_ACTIONS.GET_CONTENT: {
            this.props.setContent(payload.lines)
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