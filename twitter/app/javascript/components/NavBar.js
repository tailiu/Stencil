import React, {Component, Fragment} from "react";

import Typography from 'material-ui/Typography';

import {AppBar, Tabs, Tab} from 'material-ui';
import Toolbar from 'material-ui/Toolbar';

import Button from 'material-ui/Button';

import Menu, { MenuItem } from 'material-ui/Menu';
import Input from 'material-ui/Input';
import TextField from 'material-ui/TextField';

import Dialog, {
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
  } from 'material-ui/Dialog';

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
};

class NavBar extends Component {

    constructor(props) {
        super(props);

        this.state = {
            value : 0,
            anchorEl: null,
            tweet_box_open: false,
        }
    }

    handleTweetBoxOpen = () => {
        console.log("HERE!");
        this.setState({tweet_box_open: true });
    };

    handleTweetBoxClose = () => {
        this.setState({ tweet_box_open: false });
    };

    goToHome = e => {
		window.location = 'http://localhost:3000/pages/home';
    }

    goToMessages = e => {
		window.location = 'http://localhost:3000/pages/messages';
    }
    
    goToProfile = e => {
		window.location = 'http://localhost:3000/pages/profile';
    }

    goToNotif = e => {
		window.location = 'http://localhost:3000/pages/notifications';
    }
    
    handleChange = (event, value) => {
        this.setState({ value });
        event.preventDefault();
    };

    handleClick = e => {
        this.setState({ anchorEl: e.currentTarget });
    };

    handleClose = () => {
        this.setState({ anchorEl: null });
    };

    render() {
        return (
            <AppBar style={styles.navbar}>
                <Toolbar>
                    <Typography 
                        variant="title" 
                        style={styles.title} 
                        onClick = {this.goToHome}>
                        Twitter
                    </Typography>
                    <div style={styles.buttonGroup} >
                        <Button onClick = {this.goToHome}>Home</Button>
                        <Button onClick = {this.goToNotif}>Notifications</Button>
                        <Button onClick = {this.goToMessages}>Messages</Button>
                    </div>
                    <Input
                        placeholder="Search Twitter"
                        style={styles.input}
                    />
                    <div style={styles.profileMenuButton}>
                        <Button
                        aria-owns={this.state.anchorEl ? 'simple-menu' : null}
                        aria-haspopup="true"
                        onClick={this.handleClick}
                        >
                        User
                        </Button>
                        <Menu
                        id="simple-menu"
                        anchorEl={this.state.anchorEl}
                        open={Boolean(this.state.anchorEl)}
                        onClose={this.handleClose}
                        >
                            <MenuItem onClick = {this.goToProfile}>Profile</MenuItem>
                            <MenuItem onClick={this.handleClose}>Settings</MenuItem>
                            <MenuItem onClick={this.handleClose}>Logout</MenuItem>
                        </Menu>
                    </div>
                    <Button 
                        style={styles.tweetButton} 
                        onClick={this.handleTweetBoxOpen}>Tweet</Button>
                    <Dialog
                        open={this.state.tweet_box_open}
                        onClose={this.handleTweetBoxClose}
                        aria-labelledby="form-dialog-title"
                        >
                        <DialogTitle id="form-dialog-title">New Tweet</DialogTitle>
                        <DialogContent>
                            <DialogContentText>
                            {/* What's on your mind? */}
                            </DialogContentText>
                            <TextField
                            autoFocus
                            margin="dense"
                            id="tweet"
                            label="What's on your mind?"
                            type="email"
                            fullWidth
                            />
                        </DialogContent>
                        <DialogActions>
                            <Button onClick={this.handleTweetBoxClose} color="primary">
                            Video/Photo
                            </Button>
                            <Button onClick={this.handleTweetBoxClose} color="primary">
                            Cancel
                            </Button>
                            <Button onClick={this.handleTweetBoxClose} color="primary">
                            Tweet!
                            </Button>
                        </DialogActions>
                    </Dialog>
                </Toolbar>
            </AppBar>
        );
    }
}

export default NavBar;