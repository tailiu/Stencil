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
    }
}

class UserProfileBox extends Component{

    constructor(props){
        super(props);
        const { cookies } = this.props;
        this.state = {
            bio_box_open : false,
            user_id : cookies.get('user_id'),
            user: [],
            user_stats: [],
            avatar_symbol: '',
            user_bio: ''
        }
    }

    handleBioBoxOpen = () => {
        console.log("HERE!");
        this.setState({bio_box_open: true });
    };

    handleBioBoxClose = () => {
        this.setState({ bio_box_open: false });
    };

    componentWillMount(){

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
            //   this.goToIndex(response.data.result.user);
            }.bind(this), 1000);
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
                        "@"+this.state.user.handle + "  |  " +
                        "Followers: " + this.state.user_stats.followers + " Following: " +  this.state.user_stats.following + " Tweets: " + this.state.user_stats.tweets
                    }
                />
                <CardContent>
                    <Typography>
                        {this.state.user.bio}
                    </Typography>
                </CardContent>
                <CardActions>
                    <Button size="small" onClick={this.like}>
                        Change Photo
                    </Button>
                    <Button size="small" onClick={this.handleBioBoxOpen}>
                        Change Bio
                    </Button>
                </CardActions>
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
            </Fragment>
        );
    }
}

export default withCookies(UserProfileBox);