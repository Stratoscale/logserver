import React from 'react'
import {ProtectedRoute, __RewireAPI__ as routerUtilsRewire} from 'router/router-utils';
import {MemoryRouter} from 'react-router'
import chai, {expect} from 'chai';
import Enzyme, {mount} from 'enzyme'
import Adapter from 'enzyme-adapter-react-16';
import chaiEnzyme from 'chai-enzyme';

Enzyme.configure({ adapter: new Adapter() });

chai.use(chaiEnzyme())

const PROTECTED = 'PROTECTED'

describe('Router-Utils:', () => {

    const testRouter = ({isAuth, pathname, findService = true}) => {

      routerUtilsRewire.__Rewire__('isAuth', () => (isAuth));
      routerUtilsRewire.__Rewire__('findService', () => (findService));

      const mockComponents =
              <MemoryRouter initialEntries={[{pathname}]}>
                <ProtectedRoute path={pathname} component={(props) => <div>{PROTECTED}</div>}/>
              </MemoryRouter>

      const renderedComponent = mount(mockComponents)
      return {renderedComponent, history: renderedComponent.instance().history, wrapper: renderedComponent}
    }

    it('redirect to the login page if the user is not authenticated', () => {
      const {history} = (testRouter({isAuth: false, pathname: '/'}))
      expect(history.location.pathname).to.equal('/login');
    })

    it('saves nextPath to the router state', () => {
      const {history} = testRouter({isAuth: false, pathname: '/next-page'});
      expect(history.location.state.nextPath).to.equal('/next-page')
    })

    it('allows an authenticated users to reach their route', () => {
      const {renderedComponent} = (testRouter({isAuth: true, pathname: '/my-page'}));
      expect(renderedComponent.find('div').text()).to.be.equal(PROTECTED)
    })

    it('redirects to / if no service was found', () => {
      const {history} = (testRouter({isAuth: true, pathname: '/my-page', findService: false}));
      expect(history.location.pathname).to.be.equal('/')
    })

    it('does not redirect to /login page if current page is login (redirect-loop)', () => {
      const {history, wrapper} = testRouter({isAuth: false, pathname: '/login'})
      expect(wrapper.find('Redirect').length).to.equal(0)
      expect(history.location.pathname).to.be.equal('/login')
    })

    it('does not redirect to / if a service not found (redirect-loop)', () => {
      const {history, wrapper} = testRouter({isAuth: true, pathname: '/', findService: false})
      expect(wrapper.find('Redirect').length).to.equal(0)
      expect(history.location.pathname).to.be.equal('/')
    })


  }
)