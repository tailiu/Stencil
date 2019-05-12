import React, {Component, Fragment} from "react";
import Typography from 'material-ui/Typography';
import Card, { CardContent } from 'material-ui/Card';
import List, {
    ListItem,
    ListItemIcon,
    ListItemSecondaryAction,
    ListItemText,
  } from 'material-ui/List';
import axios from 'axios';
import { withCookies } from 'react-cookie';
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

class FollowRequestsBox extends Component{

    constructor(props){
        super(props);

        this.cookies = this.props.cookies;

        this.state = {
            bio_box_open : false,
            user_id : props.user_id,
            logged_in_user: this.cookies.get("user_id"),
            user: [],
            follow_requests : []
        }
    }

    componentDidMount(){
        this.getFollowRequests();
    }

    approveFollowRequest =(from_user, to_user)=> {

        axios.get(
            'http://localhost:8000/users/approveFollowRequest',
            {
                withCredentials: true,
                params: {
                'from_user_id': from_user,
                'to_user_id': to_user,
                "req_token": this.cookies.get('req_token')
              }
            }
          ).then(response => {
            if(response.data.result.success){
                console.log("DONE:"+from_user);
            }else{
                
            }
            this.getFollowRequests()
          })
    }

    getFollowRequests =()=> {
        // console.log("Fetch new requests!")
        axios.get(
            'http://localhost:8000/users/getFollowRequests',
            {
                withCredentials: true,
                params: {
                    'user_id': this.state.user_id, 
                    "req_token": this.cookies.get('req_token')
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
                            {this.state.follow_requests.map((req) =>
            
                                <ListItem key={req.id}>
                                    <ListItemText
                                        // primary={req.user.name + ", @"+req.user.handle}
                                        primary={
                                            renderHTML(
                                                "<strong>"+
                                                "<a style='text-decoration: none;' href='/profile/"+req.user.id+"'>"+
                                                req.user.name+
                                                "</a>"+
                                                "</strong>, <i>@"+req.user.handle+"</i>"
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
                            }
                            </List>
                        </CardContent>
                    </Card>
                }
            </Fragment>
        );
    }
}

export default withCookies(FollowRequestsBox);