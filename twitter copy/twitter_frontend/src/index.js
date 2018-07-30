import React from 'react';
import ReactDOM from 'react-dom';
import Welcome from './Components/Welcome';
import Home from './Components/Home';
import SignUp from './Components/SignUp';
import Login from './Components/Login';
import Search from './Components/Search';
import Profile from './Components/Profile';
import MessagePage from './Components/MessagePage';
import Notif from './Components/Notif';
import Settings from './Components/Settings';
import TweetPage from './Components/TweetPage';
import registerServiceWorker from './registerServiceWorker';
import {
    BrowserRouter as Router,
    Route,
    Switch,
  } from 'react-router-dom';
import { CookiesProvider } from 'react-cookie';


ReactDOM.render(
    <CookiesProvider>
    <Router>
        <Switch>
            <Route path='/home' component={Home}/>
            <Route path="/signup" component={SignUp} />
            <Route path="/login" component={Login} />
            <Route path="/search" component={Search} />
            <Route path="/profile/:user_id" component={Profile} />
            <Route path="/profile" component={Profile} />
            <Route path="/tweet/:tweet_id" component={TweetPage} />
            <Route path="/messages" component={MessagePage} />
            <Route path="/notifications" component={Notif} />
            <Route path="/settings" component={Settings} />
            <Route path="" component={Welcome} />
        </Switch>
    </Router>
    </CookiesProvider>,
    document.getElementById('root')
);

registerServiceWorker();


if (module.hot) {
    module.hot.accept();
}