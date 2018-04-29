import React, {Component} from "react";
import Avatar from 'material-ui/Avatar';
import Card, { CardHeader } from 'material-ui/Card';
import axios from 'axios';


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

    constructor(props) {

        super(props);

        this.state = {
            followers: '',
            following: '',
            tweets: ''
        };

    }

    componentDidMount() {
        this.getFollowRelationship()
    }

    getTweets = () => {
        axios.get(
            'http://localhost:3000/tweets/',
            {
                params: {
                    'id': this.props.user.id,
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

    getFollowRelationship = () => {
        axios.get(
            'http://localhost:3000/user_actions/',
            {
                params: {
                    'id': this.props.user.id,
                    "type": 'follow'
                }
            }
            ).then(response => {
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
                    title={this.props.user.name}

                    subheader={"Follwers: " + this.state.followers + " Following: " +  this.state.following}

                    // subheader="Followers:49, Following:51, Tweets:90"
                />
            </Card>
        );
    }
}

export default UserInfo;