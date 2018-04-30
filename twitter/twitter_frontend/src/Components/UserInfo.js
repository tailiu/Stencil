import React, {Component} from "react";
import Avatar from 'material-ui/Avatar';
import Card, { CardHeader } from 'material-ui/Card';
import axios from 'axios';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';

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

class UserInfo extends Component{

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    };

    constructor(props) {

        super(props);

        this.state = {
            followers: '',
            following: '',
            tweets: ''
        };

    }

    componentDidMount() {
        const { cookies } = this.props;
        
        const session_id = cookies.get("session_id")

        this.getFollowRelationship(session_id)
        // this.getTweets(session_id)
        // this.getUsernameAndHandle(session_id)
    }

    getUsernameAndHandle = (session_id) => {

    }

    getTweets = (session_id) => {
        axios.get(
            'http://localhost:3000/tweets/',
            {
                params: {
                    'id': session_id,
                    "type": 'tweet_num'
                }
            }
            ).then(response => {
                this.setState({
                    tweets: response.data.result.tweet_num,
                })

            }
        )
    }

    getFollowRelationship = (session_id) => {
        
        axios.get(
            'http://localhost:3000/user_actions/',
            {
                params: {
                    'id': session_id,
                    "type": 'follow'
                }
            }
            ).then(response => {
                console.log(JSON.stringify(response))
                this.setState({
                    followers: response.data.result.followed_num,
                    following: response.data.result.following_num
                })

            }
        )
      
    }
    

    render(){
        return(
            <Card align="left" style={styles.user_info.container}>
                <CardHeader
                    avatar={
                        <Avatar aria-label="Recipe" style={styles.user_info.avatar}>
                        TC 
                        </Avatar>
                    }
                    title={"good man"}

                    subheader={"Follwers: " + this.state.followers + " Following: " +  this.state.following + " Tweets: " + this.state.tweets}

                    // subheader="Followers:49, Following:51, Tweets:90"
                />
            </Card>
        );
    }
}

export default withCookies(UserInfo);