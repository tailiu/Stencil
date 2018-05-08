import React, {Component, Fragment} from "react";
import Button from 'material-ui/Button';
import Avatar from 'material-ui/Avatar';
import Typography from 'material-ui/Typography';
import TextField from 'material-ui/TextField';
import Card, { CardActions, CardContent, CardHeader } from 'material-ui/Card';
import List, {
    ListItem,
    ListItemAvatar,
    ListItemIcon,
    ListItemSecondaryAction,
    ListItemText,
  } from 'material-ui/List';
import axios from 'axios';
import { withCookies, Cookies } from 'react-cookie';
import MessageBar from './MessageBar';
import renderHTML from 'react-render-html';

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
    },
    follow_button: {
        backgroundColor: "#00aced",
        color: "#fff",
        marginBottom: 5
    },
    unfollow_button: {
        backgroundColor: "#F94877",
        color: "#fff",
        marginBottom: 5
    },
    pending_follow_button: {
        backgroundColor: "#ffb347",
        color: "#fff",
        marginBottom: 5
    },
    follow_list_item: {
        backgroundColor: "#e6ecf0",
        padding: 10,
        marginBottom: 10
    },
    follow_box_header: {
        // marginBottom: 20,
        fontWeight: "bold"
    },
    approve_button: {
        cursor: "pointer"
    },

}

class FollowRequestList extends Component {

    constructor(props) {
        super(props)

        this.state = {
            reqList: this.props.requests
        }
    }

    approveFollowRequest =(from_user, to_user)=> {

        axios.get(
            'http://localhost:3000/users/approveFollowRequest',
            {
              params: {
                'from_user_id': from_user,
                'to_user_id': to_user
              }
            }
          ).then(response => {
            if(response.data.result.success){
                console.log("DONE");
            }else{
                
            }
          })
    }

    render() {
        return(
            this.state.reqList.map((req) =>
            // <div style={styles.follow_list_item}>
            //     <Typography variant="body2" >
            //     {req.user.name}, @{req.user.handle}
            //     </Typography>
            // </div>
            <ListItem>
                <ListItemText
                    // primary={req.user.name + ", @"+req.user.handle}
                    primary={
                        renderHTML(
                            "<strong>"+req.user.name + "</strong>, <i>@"+req.user.handle+"</i>"
                        )
                    }
                    // secondary={secondary ? 'Secondary text' : null}
                />
                <ListItemSecondaryAction onClick={this.approveFollowRequest.bind(this, req.from_user_id, req.to_user_id)}>
                    <ListItemIcon >
                        <img style={styles.approve_button} alt="allow" src={require('../Assets/Images/approve-icon.png')} /> 
                    </ListItemIcon>
                </ListItemSecondaryAction>
            </ListItem>
            )
        )
    }

}

class FollowRequestsBox extends Component{

    constructor(props){
        super(props);
        const { cookies } = this.props;

        this.state = {
            bio_box_open : false,
            user_id : props.user_id,
            logged_in_user: cookies.get("user_id"),
            user: [],
            follow_requests : []
        }
    }

    componentDidMount(){
        this.getFollowRequests();
    }

    getFollowRequests =()=> {
        axios.get(
            'http://localhost:3000/users/getFollowRequests',
            {
                params: {
                    'user_id': this.state.user_id, 
                }
            }
          ).then(response => {
            if(response.data.result.success){
              this.setState({
                  follow_requests: response.data.result.follow_requests,
              })
            }else{
              
            }
          })
    }

    render(){
        return(
            <Fragment>
                {this.state.follow_requests.length > 0 &&
                    <Card>
                        
                        {/* <CardHeader
                            title="Follow Requests"
                        /> */}
                        <CardContent>
                            <Typography color="textSecondary" style={styles.follow_box_header}>
                                Follow Requests
                            </Typography>
                            <hr />
                            <List dense={true}>
                                <FollowRequestList requests={this.state.follow_requests} />
                            </List>
                        </CardContent>
                    </Card>
                }
            </Fragment>
        );
    }
}

export default withCookies(FollowRequestsBox);