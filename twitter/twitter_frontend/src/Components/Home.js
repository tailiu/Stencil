import React, {Component, Fragment} from "react";
import Grid from 'material-ui/Grid';
import NavBar from './NavBar';
import MessageBar from './MessageBar';
import UserInfo from './UserInfo';
import axios from 'axios';
import { instanceOf } from 'prop-types';
import { withCookies, Cookies } from 'react-cookie';
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

class Home extends Component {

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
      };

    constructor(props) {

        super(props);
        const { cookies } = this.props;

        this.state = {
            user_id: cookies.get('user_id'),
            email : '',
            name : '',
            tweet_value : '',
            value : '',
            user: '',
            tweets : []
        } 
    }

    fetchTweets =()=> {
        axios.get(
            'http://localhost:3000/tweets/mainPageTweets',
            {
                params: {
                'user_id': this.state.user_id, 
                }
            }
            ).then(response => {
                if(response.data.result.success){
                    this.setState({
                        tweets: response.data.result.tweets,
                    })
                }else{
                    // this.MessageBar.showSnackbar("User doesn't exist!");
                }
            })
    }

    componentWillMount(){
        this.fetchTweets();
        this.timer = setInterval(()=> this.fetchTweets(), 30000);
    }

    componentWillUnmount() {
        this.timer = null;
      }

    goToIndex = () => {
        this.props.history.push({pathname: '/'});
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
                    <UserInfo user={this.state.user}/>

                </Grid>

                <Grid item xs={7}>
                    <Grid container spacing={8} direction="column" align="left">
                        <TweetList tweets={this.state.tweets}/>
                    </Grid>
                </Grid>
                <Grid item xs={1}>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default withCookies(Home);
