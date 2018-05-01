import React, {Component, Fragment} from "react";
import Button from 'material-ui/Button';
import Avatar from 'material-ui/Avatar';
import Typography from 'material-ui/Typography';
import TextField from 'material-ui/TextField';
import Card, { CardActions, CardContent, CardHeader } from 'material-ui/Card';
import Dialog, {
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
  } from 'material-ui/Dialog';
  import axios from 'axios';
import { instanceOf } from 'prop-types';
import { withCookies, Cookies } from 'react-cookie';
import MessageBar from './MessageBar';
import renderHTML from 'react-render-html';

const styles = {
    user_info: {
        avatar: {

        },
        container: {
            // backgroundColor: "#00aced",
            // color: "#fff"
        },
        card: {

        }
    },
    follow_button: {
        backgroundColor: "#00aced",
        color: "#fff",
        marginBottom: 5
    },
    unfollow_button: {
        backgroundColor: "#F94877",
        color: "#fff",
        marginBottom: 5
    }
}

class UserProfileBox extends Component{

    constructor(props){
        super(props);
        const { cookies } = this.props;

        this.state = {
            bio_box_open : false,
            user_id : props.user_id,
            logged_in_user: cookies.get("user_id"),
            user: [],
            user_stats: [],
            avatar_symbol: '',
            user_bio: '',
            does_follow: false
        }
    }

    handleBioBoxOpen = () => {
        console.log("HERE!");
        this.setState({bio_box_open: true });
    };

    handleBioBoxClose = () => {
        this.setState({ bio_box_open: false });
    };

    goToHome = e => {
		window.location = '/home';
    }

    checkDoesFollow =()=> {
        axios.get(
            'http://localhost:3000/users/checkFollow',
            {
              params: {
                'from_user_id': this.state.logged_in_user, 
                'to_user_id': this.state.user_id, 
              }
            }
          ).then(response => {
            if(response.data.result.success){
              this.setState({
                  does_follow: response.data.result.follow,
              })
            }else{
              
            }
          })
    }

    getUserInfo =()=> {
        axios.get(
        'http://localhost:3000/users/getUserInfo',
        {
            params: {
            'user_id': this.state.user_id, 
            }
        }
        ).then(response => {
        if(response.data.result.success){
            this.setState({
                user: response.data.result.user,
                user_stats: response.data.result.user_stats,
                avatar_symbol: response.data.result.user.name[0]
            })
        }else{
            this.MessageBar.showSnackbar("User doesn't exist!");
            setTimeout(function() { 
            this.goToHome();
            }.bind(this), 1000);
        }
        })
    }

    componentWillMount(){
        this.getUserInfo();
        this.checkDoesFollow();
    }

    handleFollow =(follow, e)=> {
        console.log("her:", follow);
        axios.get(
            'http://localhost:3000/users/handleFollow',
            {
              params: {
                'from_user_id': this.state.logged_in_user, 
                'to_user_id': this.state.user_id, 
                'follow' : follow
              }
            }
          ).then(response => {
            if(response.data.result.success){
                this.setState({
                  does_follow: response.data.result.follow,
                })
                this.getUserInfo();
            }else{
                
            }
          })
    }

    render(){
        return(
            <Fragment>
            <Card align="left" style={styles.user_info.container}>
                <CardHeader
                    avatar={
                        <Avatar aria-label="Recipe" style={styles.user_info.avatar}>
                        {this.state.avatar_symbol}
                        </Avatar>
                    }
                    title={this.state.user.name}
                    subheader={
                        renderHTML(
                            "<strong> @"+this.state.user.handle + "</strong>  <br>  <i>" +
                            "Followers: " + this.state.user_stats.followers + ", Following: " +  this.state.user_stats.following + ", Tweets: " + this.state.user_stats.tweets + "</i>"
                        )
                    }
                />
                <CardContent>
                    <Typography>
                        {this.state.user.bio}
                    </Typography>
                </CardContent>
                {this.state.logged_in_user === this.state.user_id ?
                    <CardActions>
                        <Button size="small" onClick={this.like}>
                            Change Photo
                        </Button>
                        <Button size="small" onClick={this.handleBioBoxOpen}>
                            Change Bio
                        </Button>
                    </CardActions>
                :
                    <CardActions>
                        {this.state.does_follow ?
                            <Button onClick={this.handleFollow.bind(this, false)} fullWidth style={styles.unfollow_button} size="small" >
                                Unfollow
                            </Button>
                        :
                            <Button onClick={this.handleFollow.bind(this, true)} fullWidth style={styles.follow_button} size="small" >
                                Follow
                            </Button>
                        }
                    </CardActions>
                }
            </Card>
            <Dialog
                open={this.state.bio_box_open}
                onClose={this.handleBioBoxOpen}
                aria-labelledby="form-dialog-title"
                >
                <DialogTitle id="form-dialog-title">Change Bio</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                    {/* What's on your mind? */}
                    </DialogContentText>
                    <TextField
                    autoFocus
                    id="bio"
                    name="bio"
                    // label="What's on your mind?"
                    fullWidth
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={this.handleBioBoxClose} color="primary">
                    Cancel
                    </Button>
                    <Button onClick={this.handleBioBoxClose} color="primary">
                    Change
                    </Button>
                </DialogActions>
            </Dialog>
            <MessageBar ref={instance => { this.MessageBar = instance; }}/>
            </Fragment>
        );
    }
}

export default withCookies(UserProfileBox);