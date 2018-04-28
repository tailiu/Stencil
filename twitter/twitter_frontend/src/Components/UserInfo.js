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
            followers: 0,
            following: 0,
            tweets: 0
        };

        // console.log(JSON.stringify(props.user.id))

        // console.log(this.state.followers)

    }

    componentDidMount() {
        this.setState({
            followers: this.getFollowingRelationship()
        })
    }

    getFollowingRelationship = () => {
        axios.get(
            'http://localhost:3000/user_actions/',
            {
                params: {
                    'id': this.props.user.id,
                    "type": 'following_relationship'
                }
            }
            ).then(response => {
                console.log(response.data.result.following_num)
                console.log(response.data.result.followed_num)
                // console.log(response)
                // if(!response.data.result.success){
                //     this.showSnackbar(response.data.result.error.message)
                // }else{
                //     this.showSnackbar("Login Successful!");
                //     cookies.set('session_id', response.data.result.session_id);
                //     setTimeout(function() { 
                //         this.goToHome(response.data.result.user);
                //     }.bind(this), 1000);
                // }
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

                    subheader={this.state.followers}

                    // subheader="Followers:49, Following:51, Tweets:90"
                />
            </Card>
        );
    }
}

export default UserInfo;