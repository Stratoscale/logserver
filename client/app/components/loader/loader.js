import React from 'react'

import { Spin, Icon } from 'antd';


const Loader = ({size = 40}) => {
  const antIcon = <Icon type="loading" style={{fontSize: size}} spin/>
  return <Spin indicator={antIcon}/>
}

export default Loader