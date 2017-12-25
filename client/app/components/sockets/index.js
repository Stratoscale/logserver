import {Component} from 'react'
import {connect} from 'react-redux'
import {setFiles, socketReady} from 'sockets/socket-actions'

let socket = null

@connect(null, {socketReady, setFiles})
export default class SocketContainer extends Component {
  constructor(props) {
    super(props)
  }

  componentWillMount() {
    if (!socket) {
      socket = new WebSocket('ws://localhost:8080/ws')

      // Connection opened
      socket.addEventListener('open', (event) => {
        this.props.socketReady()
        console.log('socket opened')
      })

      // Listen for messages
      socket.addEventListener('message', (event) => {
        const {tree} = JSON.parse(event.data)
        this.props.setFiles(tree)
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