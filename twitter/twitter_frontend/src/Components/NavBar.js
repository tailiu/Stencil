import React, {Component, Fragment} from "react";
import axios from 'axios';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';

import Typography from 'material-ui/Typography';

import {AppBar, Card} from 'material-ui';
import Toolbar from 'material-ui/Toolbar';

import Button from 'material-ui/Button';

import Menu, { MenuItem } from 'material-ui/Menu';
import Input from 'material-ui/Input';
import TextField from 'material-ui/TextField';
import MessageBar from './MessageBar';
import NewTweetDialog from "./NewTweetDialog";
import { CardContent } from "material-ui";
import UserSearchBox from "./UserSearchBox";

const styles = {
    root: {
        flexGrow: 1,
      },
      flex: {
        flex: 1,
      },
    navbar: {
        height: 60,
        backgroundColor: "#fff",
        overflow: "hidden"
    },
    title: {
        color: "#00aced",
        cursor: "pointer",
    },
    tabs: {
        float: "right",
        color: "#000",
        marginRight: "0",
    },
    tab: {
        inkbar: "green", 
        underline: "#fff"
    },
    buttonGroup: {
        flex: 1,
        marginLeft: 10
    },
    tweetButton: {
        backgroundColor: "#00aced",
        color: "#fff"
    },
    profileMenuButton: {
        marginRight: 5
    },
    input: {
        marginRight: 5
    },
    notif_num: {
        // display: "inline-block",
        marginLeft: 3,
        marginTop: 3,
        fontSize: 16,
        // opacity: 0.6,
        color: "#00aced",
        fontWeight: "bold",
        fontFamily: '"Courier New", Courier, "Lucida Sans Typewriter"'
    }
};

class NavBar extends Component {

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
      };

    constructor(props) {
        
        super(props);
        
        const { cookies } = this.props;

        this.state = {
            user_id: cookies.get('user_id'),
            user_name: cookies.get('user_name'),
            user_handle: cookies.get('user_handle'),
            value : 0,
            anchorEl: null,
            tweet_box_open: false,
            search_query: "",
            notifications: 0,
        }

    }

    componentDidMount() {
        this.getNotifications();
        this.timer = setInterval(()=> this.getNotifications(), 30000);
    }

    componentWillUnmount() {
        this.timer = null;
    }

    handleTweetBoxClose = () => {
        this.setState({ tweet_box_open: false });
    };

    handleTweetBoxOpen = () => {
        console.log("HERE!");
        this.setState({tweet_box_open: true });
    };

    goToIndex = e => {
		window.location = '/';
    }

    goToHome = e => {
		window.location = '/home';
    }

    goToMessages = e => {
		window.location = '/messages';
    }
    
    goToProfile = e => {
		window.location = '/profile';
    }

    goToSettings = e => {
		window.location = '/settings';
    }

    goToNotif = e => {
		window.location = '/notifications';
    }

    handleClick = e => {
        this.setState({ anchorEl: e.currentTarget });
    };

    handleClose = () => {
        this.setState({ anchorEl: null });
    };

    handleLogout = () =>  {
        const { cookies } = this.props;
        axios.get(
        'http://localhost:3000/users/logout'
        ).then(response => {
            cookies.remove('session_id');
            this.goToIndex();
        })
      }

    getNotifications = () =>  {

        axios.get(
        'http://localhost:3000/notifications/getNewNotifications',
        {
            params: {
              'user_id': this.state.user_id, 
            }
          }
        ).then(response => {
            console.log(response)
            if(response.data.result.success){
                this.setState({
                  notifications: response.data.result.notifs,
                })

            }else{
                // console.log("Notif error!");
            }
          })
      }

    searchUser =(e)=> {
        // console.log(e.target.value)
        this.setState({
            search_query: e.target.value
        })
    }

    render() {
        return (
            <Fragment>
            <AppBar style={styles.navbar}>
                <MessageBar ref={instance => { this.MessageBar = instance; }}/>
                <Toolbar>
                    <Typography 
                        variant="title" 
                        style={styles.title} 
                        onClick = {this.goToHome}>
                        Twitter
                    </Typography>
                    <div style={styles.buttonGroup} >
                        <Button onClick = {this.goToHome}>Home</Button>
                        <Button onClick = {this.goToNotif}>
                            Notifications
                            <Typography 
                                variant="button"
                                style={styles.notif_num}>
                                {this.state.notifications}
                            </Typography>
                        </Button>
                        <Button onClick = {this.goToMessages}>Messages</Button>
                    </div>
                    <Input
                        placeholder="Search Twitter"
                        style={styles.input}
                        onChange={this.searchUser}
                    />
                    <div style={styles.profileMenuButton}>
                        <Button
                        aria-owns={this.state.anchorEl ? 'simple-menu' : null}
                        aria-haspopup="true"
                        onClick={this.handleClick}
                        >
                        {this.state.user_name}
                        </Button>
                        <Menu
                        id="simple-menu"
                        anchorEl={this.state.anchorEl}
                        open={Boolean(this.state.anchorEl)}
                        onClose={this.handleClose}
                        >
                            <MenuItem onClick = {this.goToProfile}>Profile</MenuItem>
                            <MenuItem onClick={this.goToSettings}>Settings</MenuItem>
                            <MenuItem onClick={this.handleLogout}>Logout</MenuItem>
                        </Menu>
                    </div>
                    <Button 
                        style={styles.tweetButton} 
                        onClick={this.handleTweetBoxOpen}>
                        Tweet
                    </Button>
                    <NewTweetDialog open={this.state.tweet_box_open} onChange={this.handleTweetBoxClose} />
                </Toolbar>
            </AppBar>

                <UserSearchBox query={this.state.search_query}/>

            </Fragment>
        );
    }
}

export default withCookies(NavBar);