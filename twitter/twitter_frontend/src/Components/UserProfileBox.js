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
import { withCookies } from 'react-cookie';
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
    },
    pending_follow_button: {
        backgroundColor: "#ffb347",
        color: "#fff",
        marginBottom: 5
    }
}

class UserProfileBox extends Component{

    constructor(props){
        super(props);
        this.cookies = this.props.cookies;

        this.state = {
            bio_box_open : false,
            user_id : props.user_id,
            logged_in_user: this.cookies.get("user_id"),
            user: [],
            user_stats: [],
            avatar_symbol: '',
            user_bio: '',
            new_bio: '',
            does_follow: false,
            does_block: false,
            does_mute: false
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

    checkDoesBlock =()=> {
        axios.get(
            'http://localhost:3000/users/checkBlock',
            {
              params: {
                'from_user_id': this.state.logged_in_user, 
                'to_user_id': this.state.user_id, 
              }
            }
          ).then(response => {
            if(response.data.result.success){
              this.setState({
                  does_block: response.data.result.block,
              })
            }else{
              
            }
          })
    }

    checkDoesMute =()=> {
        axios.get(
            'http://localhost:3000/users/checkMute',
            {
              params: {
                'from_user_id': this.state.logged_in_user, 
                'to_user_id': this.state.user_id, 
              }
            }
          ).then(response => {
            if(response.data.result.success){
              this.setState({
                  does_mute: response.data.result.mute,
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
            this.checkDoesFollow();
            this.checkDoesBlock();
            this.checkDoesMute();
        }else{
            this.MessageBar.showSnackbar("User doesn't exist!");
        }
        })
    }

    componentWillMount(){
        this.getUserInfo();
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
                this.getUserInfo();
                console.log("followd: success");
                console.log(response);
            }else{
                console.log("followd: failed");
                console.log(response);
                this.MessageBar.showSnackbar(response.data.result.error.message);   
            }
          })
    }

    handleBlock =(block, e)=> {
        console.log("her:", block);
        axios.get(
            'http://localhost:3000/users/handleBlock',
            {
                params: {
                    'from_user_id': this.state.logged_in_user, 
                    'to_user_id': this.state.user_id, 
                    'block' : block
                }
            }
        ).then(response => {
            if(response.data.result.success){
                this.setState({
                  does_block: response.data.result.block,
                })
                this.getUserInfo();
            }else{
                console.log("block: failed");
                console.log(response);
                this.MessageBar.showSnackbar(response.data.result.error.message);   
            }
        })

        axios.get(
            'http://localhost:3000/conversations/blockInGroupConversation',
            {
                params: {
                    'from_user_id': this.state.logged_in_user, 
                    'to_user_id': this.state.user_id
                }
            }
        ).then(response => {
            if(response.data.result.success){
            }else{
                this.MessageBar.showSnackbar(response.data.result.error.message);   
            }
        })
    }

    handleMute =(mute, e)=> {
        console.log("her:", mute);
        axios.get(
            'http://localhost:3000/users/handleMute',
            {
              params: {
                'from_user_id': this.state.logged_in_user, 
                'to_user_id': this.state.user_id, 
                'mute' : mute
              }
            }
          ).then(response => {
            if(response.data.result.success){
                this.setState({
                  does_mute: response.data.result.mute,
                })
                this.getUserInfo();
            }else{
                console.log("mute: failed");
                console.log(response);
                this.MessageBar.showSnackbar(response.data.result.error.message);      
            }
          })
    }

    handleBio =(e)=> {
        console.log("her:", this.state.new_bio);
        axios.get(
            'http://localhost:3000/users/updateBio',
            {
              params: {
                'user_id': this.state.logged_in_user, 
                'bio': this.state.new_bio
              }
            }
          ).then(response => {
            if(response.data.result.success){
                this.setState({
                  user: response.data.result.user,
                })
                this.handleBioBoxClose();
                this.MessageBar.showSnackbar("Bio updated!");
            }else{
                this.MessageBar.showSnackbar("Bio can't be updated!");
            }
          })
    }

    changeBioContent = (e) => {
        this.setState({
            new_bio: e.target.value
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
                        {/* <Button size="small" onClick={this.like}>
                            Change Photo
                        </Button> */}
                        <Button size="small" onClick={this.handleBioBoxOpen}>
                            Change Bio
                        </Button>
                    </CardActions>
                :
                    <CardActions>
                        {!this.state.does_block ?
                            this.state.does_follow === "pending" ?
                                <Button onClick={this.handleFollow.bind(this, true)} fullWidth style={styles.pending_follow_button} size="small" >
                                    Pending
                                </Button>
                            :   
                            
                                this.state.does_follow ?
                                    <Button onClick={this.handleFollow.bind(this, false)} fullWidth style={styles.unfollow_button} size="small" >
                                        Unfollow
                                    </Button>
                                :
                                    <Button onClick={this.handleFollow.bind(this, true)} fullWidth style={styles.follow_button} size="small" >
                                        Follow
                                    </Button>
                            :
                            <Fragment>
                            </Fragment>
                                
                        }
                        {this.state.does_block ?
                            <Button onClick={this.handleBlock.bind(this, false)} fullWidth style={styles.unfollow_button} size="small" >
                                Unblock
                            </Button>
                        :
                            <Button onClick={this.handleBlock.bind(this, true)} fullWidth style={styles.follow_button} size="small" >
                                Block
                            </Button>
                        }
                        {!this.state.does_block?
                            this.state.does_mute ?
                                <Button onClick={this.handleMute.bind(this, false)} fullWidth style={styles.unfollow_button} size="small" >
                                    Unmute
                                </Button>
                            :
                                <Button onClick={this.handleMute.bind(this, true)} fullWidth style={styles.follow_button} size="small" >
                                    Mute
                                </Button>
                            :
                            <Fragment>
                            </Fragment>
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
                    onChange={this.changeBioContent}
                    fullWidth
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={this.handleBioBoxClose} color="primary">
                    Cancel
                    </Button>
                    <Button onClick={this.handleBio.bind(this)} color="primary">
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