import React, {Component, Fragment} from "react";
import Grid from 'material-ui/Grid';
import Card, { CardContent, CardHeader } from 'material-ui/Card';
import NavBar from './NavBar';
import List, { ListItem, ListItemIcon, ListItemText } from 'material-ui/List';
import axios from 'axios';
import { withCookies } from 'react-cookie';
import Typography from 'material-ui/Typography';
import renderHTML from 'react-render-html';
import Moment from 'moment';

const styles = {
    grid : {
        container : {
            marginTop: 80,
            height: 200
        }
    },
    new_notif: {
        backgroundColor: "#c0deed",
        marginBottom: 5
    },
    old_notif: {
        marginBottom: 5
    }
};

class Notif extends Component {

    constructor(props) {

        super(props);
        
        this.cookies = this.props.cookies;
        
        this.state = {
            user_id: this.cookies.get('user_id'),
            notifs: []
        }
    }

    componentDidMount(){
        // console.log("here notif")
        this.getNotifs();
    }

    getNotifs =(e)=> {

        axios.get(
            'http://localhost:3000/notifications/get',
            {
              params: {
                'user_id': this.state.user_id,
                "req_token": this.cookies.get('req_token') 
              }
            }
          ).then(response => {
            // console.log(response)
            if(response.data.result.success){
                this.setState({
                  notifs: response.data.result.notifs,
                })

            }else{
                console.log("Notif error!");
            }
          })
    }

    getNotifStyle = (isSeen) => {
        if (isSeen) return styles.old_notif
        else return styles.new_notif
    }

    render () {
        const notif_info = {
            'retweet': {
                'text': "retweeted",
                'icon': require('../Assets/Images/retweet_notif.png')
            },
            'like': {
                'text': "liked",
                'icon': require('../Assets/Images/like_notif.png')
            },
            'reply': {
                'text': "replied to",
                'icon': require('../Assets/Images/reply_icon.png')
            },
            'follow': {
                'text': "followed you",
                'icon': require('../Assets/Images/follow_notif.png')
            },
        }
    return (
        <Fragment>
            <NavBar />
            <Grid style={styles.grid.container} direction="column" container spacing={24} >
                
                <Grid item xs={2}>
                </Grid>
                <Grid item xs={8}>
                    <Card>
                        <CardHeader
                            title="Notifications"
                        />
                        <hr />
                        <CardContent>
                        <List dense={this.state.notifs.length > 0}>
                            {this.state.notifs.length <= 0? 
                                <ListItem>
                                    <ListItemText primary='No Notifications!' />
                                </ListItem>
                            :

                                <Fragment>
                                {this.state.notifs.map((notif) =>
                                    
                                    <ListItem style={this.getNotifStyle(notif.is_seen)}>
                                        
                                        <ListItemIcon>
                                        <img style={styles.logo} alt="retweet" src={notif_info[notif.notification_type]["icon"]} /> 
                                        </ListItemIcon>
                                        {notif.notification_type == "follow"? 
                                            <ListItemText
                                            primary={
                                                renderHTML("<a style='text-decoration: none;' href='/profile/"+notif.from_user_id+"'>"+notif.from_user_name+"</a>" + " " + notif_info[notif.notification_type]["text"])
                                            }
                                            secondary={
                                                Moment(notif.created_at).format('MMMM Do, YYYY - h:mm A')
                                            }
                                            />
                                        :
                                            <ListItemText
                                            primary={
                                                renderHTML("<a style='text-decoration: none;' href='/profile/"+notif.from_user_id+"'>"+notif.from_user_name+"</a>"  + " " + notif_info[notif.notification_type]["text"]+" your <a style='text-decoration: none;' href='/tweet/"+notif.tweet+"'>tweet</a>")
                                            }
                                            secondary={
                                                Moment(notif.created_at).format('MMMM Do, YYYY - h:mm A')
                                            }
                                            />
                                        
                                        }
                                    </ListItem>
                                )}
                                </Fragment>

                            }
                            </List>
                        </CardContent>
                    </Card>
                </Grid>
                <Grid item xs={2}>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default withCookies(Notif);
