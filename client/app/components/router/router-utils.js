import React from 'react'
import {Route, Redirect} from 'react-router-dom'

const RedirectToLogin = (props) => {
  console.info('user is not authenticated, redirecting to the login page');
  const isLoginPage = props.location.pathname === '/login'
  if (__DEV__) {
    if (isLoginPage) {
      console.error('ProtectedRoute was defined on /login, it should be a normal Route')
    }
  }

  return !isLoginPage ? <Redirect to={{
    pathname: '/login', state: {
      nextPath: props.location.pathname,
      search:   props.location.search,
    },
  }}/> : null
}

const wrapComponentWithDevCheck = (component, routeProps) => (props) => {
  const baseURL        = props.location.pathname.split('/')[1]
  const isServiceFound = props.location.pathname === '/' || findService(`/${baseURL}`, 'url') || findPage(baseURL)

  if (__DEV__) {
    console.info(`route "${props.match.path}" found for url "${props.location.pathname}" with component:`, component && component.name.toLowerCase() === 'connect' ? `Connected(${component.WrappedComponent.name})` : component.name)
    if (!isServiceFound) {
      console.error(`route for service ${baseURL} was not found,
        this might be an authorization issue, or the service was not defined in services.js`)
    }
  }

  if (isServiceFound && routeProps.title && !props.location.search) {
    if (routeProps.path.includes(':entity_id')) {
      track(`${isServiceFound.title}-${routeProps.title}-Entity`)
    }
    else {
      track(`${isServiceFound.title}-${routeProps.title}`)
    }
  }

  return isServiceFound ? React.createElement(component, {...props, ...routeProps}) : <Redirect to={{pathname: '/'}}/>
}


export const ProtectedRoute = ({component, ...rest}) => {
  if (!component) {
    throw new Error('can not create a route without a component')
  }

  return (<Route {...rest} render={isAuth() ?
    wrapComponentWithDevCheck(component, rest) : RedirectToLogin}/>)
}

export function generateRoutes(serviceUrl, {url, component, entityComponent, subItems, title}) {
  const routes = [];
  const key    = serviceUrl;
  if (subItems) {
    const defaultSubItem      = subItems[0];
    const defaultSubItemRoute = generateRoutes(url, defaultSubItem);
    return subItems.map(item => generateRoutes(`${serviceUrl}/${item.url}`, item)).concat([defaultSubItemRoute]);
  } else {
    if (entityComponent) {
      routes.push({path: `${serviceUrl}/:entity_id`, component: entityComponent, key: `${serviceUrl}-single`, title})
    }

    if (component) {
      routes.push({path: `${serviceUrl}`, component, key, title});
    }
  }

  return routes.map(({path, component, key}) => <ProtectedRoute key={serviceUrl} title={title} path={path} component={component}/>)
}