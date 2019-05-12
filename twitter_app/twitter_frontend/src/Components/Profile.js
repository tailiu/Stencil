import React, {Component, Fragment} from "react";
import Grid from 'material-ui/Grid';
import NavBar from './NavBar';
import UserProfileBox from './UserProfileBox';
import axios from 'axios';
import { withCookies } from 'react-cookie';
import MessageBar from './MessageBar';
import TweetList from './TweetList';
import FollowRequestsBox from "./FollowRequestsBox";

const styles = {
    grid : {
        container : {
            marginTop: 80
        }
    },
    card: {
        card:{
            // minWidth: 400,
        },
        input:{
            width: "95%",
        },
        button: {
            width: "100%",
            backgroundColor: "#00aced",
            color: "#fff",
        }
    },
    tweet: {
        avatar: {

        },
        main_input: {

        },
        container: {

        },
        card: {

        }
    }
};

class Profile extends Component {

    constructor(props) {

        super(props);

        this.cookies = this.props.cookies;

        let user_id = this.cookies.get('user_id');

        if (props.match.params.user_id){
            user_id = props.match.params.user_id
        }

        this.state = {
            user_id: user_id,
            logged_in_user: this.cookies.get('user_id'),
            session_id: this.cookies.get('session_id'),
            email : '',
            password : '',
            name : '',
            tweet_value : '',
            value : '',
            tweets: []
        }

    }

    componentWillMount(){

        axios.get(
            'http://localhost:8000/users/checkTwoWayBlock',
            {
                withCredentials: true,
                params: {
                'from_user_id': this.state.logged_in_user, 
                'to_user_id': this.state.user_id, 
                "req_token": this.cookies.get('req_token')
              }
            }
          ).then(response => {
            
            if(response.data.result.success){
                if(response.data.result.block){
                    this.MessageBar.showSnackbar("BLOCKED");
                }else{
                    this.fetchTweets();
                    this.timer = setInterval(()=> this.fetchTweets(), 30000);
                }
            }else{
              
            }
          })
    }

    componentWillUnmount() {
        clearInterval(this.timer);
      }

    fetchTweets =()=> {

        axios.get(
        'http://localhost:8000/tweets/fetchUserTweets',
        {
            withCredentials: true,
            params: {
            'user_id': this.state.user_id, 
            'requesting_user': this.state.logged_in_user,
            'session_id': this.state.session_id,
            "req_token": this.cookies.get('req_token')
            }
        }
        ).then(response => {

            if(response.data.result.success){
                // console.log("result:")
                // console.log(response.data)
                this.setState({
                    tweets: response.data.result.tweets,
                })
            }else{
                this.MessageBar.showSnackbar(response.data.result.error.message);
                clearInterval(this.timer);
            }
        })
    }

    handleChange =(e)=> {
        this.setState({ value: e.target.value });
    }

    handleChangeToTweetField =(e)=> {
        this.setState({ tweet_value: e.target.value });
    }

    render () {
    return (
        <Fragment>
            <NavBar />
            <MessageBar ref={instance => { this.MessageBar = instance; }}/>
            <Grid style={styles.grid.container} container spacing={16}>
                
                <Grid item xs={1}>
                </Grid>

                <Grid item xs={3}>
                    <Grid container direction="column" spacing={8}>
                        <Grid item>
                            <UserProfileBox user_id={this.state.user_id}/>
                        </Grid>
                        <Grid item>
                            {this.state.logged_in_user === this.state.user_id &&
                                <FollowRequestsBox user_id={this.state.user_id}/>
                            }
                        </Grid>
                    </Grid>
                </Grid>

                <Grid item xs={7}>
                    <Grid container spacing={8} direction="column" align="left">
                        <TweetList tweets={this.state.tweets} />
                    </Grid>
                </Grid>
                <Grid item xs={1}>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default withCookies(Profile);
