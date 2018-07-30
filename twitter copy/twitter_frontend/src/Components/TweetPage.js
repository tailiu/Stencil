import React, {Component, Fragment} from "react";
import Grid from 'material-ui/Grid';
import NavBar from './NavBar';

import axios from 'axios';
import { withCookies } from 'react-cookie';
import MessageBar from './MessageBar';
import TweetList from './TweetList';

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
        
        this.state = {
            tweets: [],
            tweet_id: props.match.params.tweet_id
        }

    }

    componentWillMount(){
        this.fetchTweets();
        this.timer = setInterval(()=> this.fetchTweets(), 30000);
    }

    componentWillUnmount() {
        this.timer = null;
      }

    fetchTweets =()=> {

        axios.get(
        'http://localhost:3000/tweets/getTweet',
        {
            withCredentials: true,
            params: {
            'tweet_id': this.state.tweet_id, 
            "req_token": this.cookies.get('req_token')
            }
        }
        ).then(response => {
            if(response.data.result.success){
                this.setState({
                    tweets: response.data.result.replies,
                })
            }else{
                
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
                <Grid item xs={2}>
                </Grid>
                <Grid item xs={8}>
                    <Grid container spacing={8} direction="column" align="left">
                        <TweetList tweets={this.state.tweets} />
                    </Grid>
                </Grid>
                <Grid item xs={2}>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default withCookies(Profile);
