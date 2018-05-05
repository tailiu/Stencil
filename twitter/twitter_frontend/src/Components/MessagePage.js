import React, {Component, Fragment} from "react";
import TextField from 'material-ui/TextField';
import Grid from 'material-ui/Grid';
import Divider from 'material-ui/Divider';
import Card, { CardContent, CardHeader } from 'material-ui/Card';
import NavBar from './NavBar';
import Avatar from 'material-ui/Avatar';
import List, { ListItem, ListItemText, } from 'material-ui/List';
import Button from 'material-ui/Button';
import Dialog, {
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
} from 'material-ui/Dialog';
import MessageBar from './MessageBar';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
import axios from 'axios';
import ConversationList from './ConversationList'
import MessageList from './MessageList'
import MessageInput from './MessageInput'

const styles = {
    grid : {
        container: {
            marginTop: 80
        }
    },
    
    card: {
        height: "500px",
        overflowY: 'auto',
        tweetButton: {
            backgroundColor: "#00aced",
            color: "#fff"
        },
    }
};



class MessagePage extends Component {

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    };

    constructor(props) {

        super(props);

        const { cookies } = this.props;

        this.state = {
            user_id : cookies.get('user_id'),
            user_handle: cookies.get('user_handle'),
            new_message_box_open: false,
            message_to: '',
            conversations: [],
            current_conversation_id: '',
            messages: ''
        }
    }

    componentDidMount() {
        this.getConversationList()
    }

    getConversationList = callback => {
        axios.get(
            'http://localhost:3000/conversations/',
            {
                params: {
                    'id': this.state.user_id
                }
            }
        ).then(response => {
            if(!response.data.result.success){
                this.MessageBar.showSnackbar(response.data.result.error.message)
            }else{
                this.setState({
                    'conversations': response.data.result.conversations
                });
                if (callback) callback()
            }
        })
    }

    getMessageList = (current_conversation_id) => {
        axios.get(
            'http://localhost:3000/messages',
            {
                params: {
                    "id": current_conversation_id
                }
            }
        ).then(response => {
            if(!response.data.result.success){
            }else{
                if (response.data.result.messages == undefined) {
                    this.setState({
                        messages: ""
                    })
                } else {
                    this.setState({
                        messages: response.data.result.messages
                    })
                }
    
            }
        })
    }

    handleNewMessageBoxOpen = e => {
        this.setState({new_message_box_open: true });
    }

    handleNewMessageBoxClose = e => {
        this.setState({new_message_box_open: false });
    }

    updateMessageTo = e => {
        this.setState({
            message_to: e.target.value
        })
    }

    validateInput = e => {
        return true
    }

    handleNewConversation = e => {
        if(!this.validateInput()){
            this.MessageBar.showSnackbar("Please input valid user handles")
            return
        }

        const raw_data = this.state.message_to.split('@')
        raw_data.shift()

        const participants = []
        for (var i in raw_data) {
            raw_data[i] = raw_data[i].replace(/\s/g,''); // replace all spaces in handles
            participants.push(raw_data[i])
        }
        participants.push(this.state.user_handle)

        axios.get(
            'http://localhost:3000/conversations/new',
            {
                params: {
                    'participants': participants
                }
            }
        ).then(response => {
            if(!response.data.result.success){
                this.MessageBar.showSnackbar(response.data.result.error.message)
            }else{
                this.getConversationList(this.handleNewMessageBoxClose)
            }
        })

    }

    handleConversationChange = current_conversation_id => {
        this.setState({
            current_conversation_id: current_conversation_id
        })
        this.getMessageList(current_conversation_id)    
    }

    handleNewMessage = () => {
        this.getMessageList(this.state.current_conversation_id)    
    }

    render () {
        return (
            <Fragment>
                <MessageBar ref={instance => { this.MessageBar = instance; }}/>
                <NavBar />
                <Grid style={styles.grid.container} container spacing={24} >
                    
                    <Grid item xs={1}>
                    </Grid>
                    <Grid item xs={10} >
                        <Grid direction="column" align="left">
                            <Grid item>
                            <Card style={styles.card}>
                                <CardHeader
                                    title="Messages"

                                    action={
                                        <Button 
                                        style={styles.card.tweetButton} 
                                        onClick={this.handleNewMessageBoxOpen}>New Message</Button>
                                    }
                                />

                                <Dialog
                                    open={this.state.new_message_box_open}
                                    onClose={this.handleNewMessageBoxClose}
                                    aria-labelledby="form-dialog-title"
                                    >
                                    <DialogTitle id="form-dialog-title">New Message</DialogTitle>
                                    <DialogContent>
                                        <DialogContentText>
                                        {/* What's on your mind? */}
                                        </DialogContentText>
                                        <TextField
                                            autoFocus
                                            margin="dense"
                                            id="tweet"
                                            label="Send message to"
                                            type="email"
                                            value={this.state.message_to}
                                            onChange={this.updateMessageTo}
                                            fullWidth
                                        />
                                    </DialogContent>
                                    <DialogActions>
                                        <Button onClick={this.handleNewMessageBoxClose} color="primary">
                                            Cancel
                                        </Button>
                                        <Button onClick={this.handleNewConversation} color="primary">
                                            New Message
                                        </Button>
                                    </DialogActions>
                                </Dialog>


                                <hr />

                                <CardContent>
                                    <Grid container direction="row" spacing={8} align="left">
                                        <Grid item xs={4}>
                                            <ConversationList 
                                                conversations = {this.state.conversations}
                                                onConversationChange =  {this.handleConversationChange}
                                            />
                                        </Grid>
                                        <Grid item xs={8} >
                                            <MessageList 
                                                messages = {this.state.messages} 
                                            />

                                            <MessageInput 
                                                current_conversation_id = {this.state.current_conversation_id}
                                                onNewMessage = {this.handleNewMessage}
                                            />
                                        </Grid>
                                    </Grid>                                
                                </CardContent>
                                </Card>
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item xs={1}>
                    </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default withCookies(MessagePage);
