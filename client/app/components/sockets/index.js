import {Component} from 'react'
import {connect} from 'react-redux'

let socket = null

@connect()
export default class SocketContainer extends Component {
  constructor(props) {
    super(props)
  }

  componentWillMount() {
    if (!socket) {
      socket = new WebSocket('ws://localhost:8080/ws')

      // Connection opened
      socket.addEventListener('open', function (event) {
        socket.send(JSON.stringify({
          action: 'get-file-tree',
        }))
      })

      // Listen for messages
      socket.addEventListener('message', function (event) {
        console.log('Message from server ', event.data)
      })

    }
  }

  render = () => {
    return null
  }
}