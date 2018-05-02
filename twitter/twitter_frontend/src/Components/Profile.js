import React, {Component, Fragment} from "react";
import Grid from 'material-ui/Grid';
import NavBar from './NavBar';
import UserProfileBox from './UserProfileBox';
import axios from 'axios';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
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

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    };

    constructor(props) {

        super(props);

        const { cookies } = this.props;

        let user_id = cookies.get('user_id');

        if (props.match.params.user_id){
            user_id = props.match.params.user_id
        }

        this.state = {
            user_id: user_id,
            email : '',
            password : '',
            name : '',
            tweet_value : '',
            value : '',
            tweets: []
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
        'http://localhost:3000/tweets/fetchUserTweets',
        {
            params: {
            'user_id': this.state.user_id, 
            }
        }
        ).then(response => {

            if(response.data.result.success){
                console.log("result:")
                console.log(response.data.result.tweets)
                this.setState({
                    tweets: response.data.result.tweets,
                })
            }else{
                this.MessageBar.showSnackbar("User doesn't exist!");
                setTimeout(function() { 
                //   this.goToIndex(response.data.result.user);
                }.bind(this), 1000);
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
                    <UserProfileBox user_id={this.state.user_id}/>

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
