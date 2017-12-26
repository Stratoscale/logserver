import React from 'react'
import SocketContainer from 'sockets'
import Home from 'layouts/home'

const App = (props) => {
  return (
    <div>
      <SocketContainer/>
      <Home/>
    </div>
  )
}

export default App