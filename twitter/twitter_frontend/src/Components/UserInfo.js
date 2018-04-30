import React, {Component} from "react";
import Avatar from 'material-ui/Avatar';
import Card, { CardHeader } from 'material-ui/Card';
import axios from 'axios';
import { instanceOf } from 'prop-types';
import { withCookies, Cookies } from 'react-cookie';
import MessageBar from './MessageBar';

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
        const { cookies } = this.props;
        this.state = {
            user_id : cookies.get('user_id'),
            user: [],
            user_stats: [],
            avatar_symbol: ''
        };

    }

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
                    // subheader={"Followers: " + this.state.user_stats.followed + " Following: " +  this.state.user_stats.following + " Tweets: " + this.state.user_stats.tweets}
                    // subheader="Followers:49, Following:51, Tweets:90"
                />
                <MessageBar ref={instance => { this.MessageBar = instance; }}/>
            </Card>
        );
    }
}

export default withCookies(UserInfo);