import React from 'react';
import axios from 'axios';
import ReactDOM from 'react-dom';
import Welcome from './Components/Welcome';
import Home from './Components/Home';
import SignUp from './Components/SignUp';
import Login from './Components/Login';
import Search from './Components/Search';
import Profile from './Components/Profile';
import Messages from './Components/Messages';
import Notif from './Components/Notif';
import Settings from './Components/Settings';
import registerServiceWorker from './registerServiceWorker';
import {
    BrowserRouter as Router,
    Route,
    Switch,
    Redirect
  } from 'react-router-dom';
import { CookiesProvider } from 'react-cookie';


function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return false;
}

function isLoggedIn(){
    // return true;
    const session_id = getCookie("session_id");
    if (session_id)
        return true;
    else return false;
}

ReactDOM.render(
    <CookiesProvider>
    <Router>
        <Switch>
            <Route
                path='/home'
                render={(props) => <Home {...props} />}
            />
            <Route path="/signup" render={() => (
                isLoggedIn() ? (
                    <Home />
                ) : (
                    <SignUp />
                )
                )}
            />
            <Route path="/login"  render={() => (
                isLoggedIn() ? (
                    <Home />
                ) : (
                    <Login />
                )
                )}
            />
            <Route path="/search"  render={() => (
                isLoggedIn() ? (
                    <Search />
                ) : (
                    <Welcome />
                )
                )}
            />
            <Route path="/profile"  render={() => (
                isLoggedIn() ? (
                    <Profile />
                ) : (
                    <Welcome />
                )
                )}
            />
            <Route path="/messages"  render={() => (
                isLoggedIn() ? (
                    <Messages />
                ) : (
                    <Welcome />
                )
                )}
            />
            <Route path="/notifications"  render={() => (
                isLoggedIn() ? (
                    <Notif />
                ) : (
                    <Welcome />
                )
                )}
            />
            <Route path="/settings"  render={() => (
                isLoggedIn() ? (
                    <Settings />
                ) : (
                    <Welcome />
                )
                )}
            />
            <Route path=""  render={() => (
                isLoggedIn() ? (
                    <Home />
                ) : (
                    <Welcome />
                )
                )}
            />
        </Switch>
  </Router>
  </CookiesProvider>,
    document.getElementById('root')
);

registerServiceWorker();


if (module.hot) {
    module.hot.accept();
  }