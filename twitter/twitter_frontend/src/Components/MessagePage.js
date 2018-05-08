import React, {Component} from "react";
import Grid from 'material-ui/Grid';
import Divider from 'material-ui/Divider';
import NavBar from './NavBar';
import Button from 'material-ui/Button';
import MessageBar from './MessageBar';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
import axios from 'axios';
import ConversationList from './ConversationList';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import Paper from 'material-ui/Paper';
import Typography from 'material-ui/Typography';
import NewConversation from './NewConversation'


const styles = {
    headerContainer: {
        marginTop: 80
    },
    header: {
        padding: 20,
        height: 35
    },
    headline: {
        float: "left"
    },
    newMessageButton: {
        backgroundColor: "#00aced",
        color: "#fff",
        variant: "raised",
        display: "inline-block",
        float: "right"
    },
    conversationListContainer: {
        height: "70vh",
        overflow: "auto"
    },
    messageListContainer: {
        height: "70vh",
        overflow: "auto"
    },
    messageList: {
        height: "58vh",
        overflow: "auto"
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
            conversations: [],
            current_conversation_id: '',
            current_conversation_type: '',
            messages: ''
        }
    }

    componentDidMount() {
        this.initialize()
        this.timer = setInterval(()=> this.periodicActions(), 6000);
    }

    componentWillUnmount() {
        this.timer = null;
    }

    periodicActions = () => {
        this.getConversationList()
        this.getMessageList(this.state.current_conversation_id)
    }

    setMessageState = (messages) => {
        this.setState({
            'messages': messages,
        });
    }

    initialize = () => {
        this.getConversationList((conversations) => {
            if (conversations.length >= 1) {
                const conversation_id = conversations[0].conversation.id 
                const conversation_type = conversations[0].conversation_type

                this.setCurrentConversation(conversation_id, conversation_type)
                this.getMessageList(conversation_id)
            } else {
                this.setMessageState('')
            }
        })
    }

    getConversationList = (cb) => {
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
                const conversations = response.data.result.conversations
                
                this.setState({
                    'conversations': conversations,
                });

                if (cb) cb(conversations)
            }
        })
    }

    setCurrentConversation = (current_conversation_id, current_conversation_type) => {
        this.setState({
            'current_conversation_id': current_conversation_id,
            'current_conversation_type': current_conversation_type
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
                    this.setMessageState('')
                } else {
                    this.setMessageState(response.data.result.messages)
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

    handleNewConversation = (new_conversation_ID, new_conversation_type) => {
        this.getConversationList()
        this.setCurrentConversation(new_conversation_ID, new_conversation_type)
        this.getMessageList(new_conversation_ID)
    }

    handleConversationChange = (current_conversation_id, current_conversation_type) => {
        this.setCurrentConversation(current_conversation_id, current_conversation_type)
        this.getMessageList(current_conversation_id)    
    }

    handleLeaveConversation = () => {
        this.initialize()
    }

    handleNewMessage = () => {
        this.getMessageList(this.state.current_conversation_id)    
    }

    render () {
        return (
            <div>
                <MessageBar ref={instance => { this.MessageBar = instance; }}/>

                <NavBar />

                <Grid style={styles.headerContainer} container spacing={24} >
                    <Grid item xs={1}>
                    </Grid>
                    <Grid item xs={10} >
                        <Paper style={styles.header}>
                            <Typography  style={styles.headline} variant="headline" component="h4">
                                Messages
                            </Typography>
                            <Button  style={styles.newMessageButton} onClick={this.handleNewMessageBoxOpen}>
                                New Message
                            </Button>
                        </Paper>
                    </Grid>
                    <Grid item xs={1}>
                    </Grid>
                </Grid>

                <Grid container spacing={24} >
                    <Grid item xs={1}>
                    </Grid>
                    <Grid item xs={3}>
                        <Paper style={styles.conversationListContainer} >
                            <ConversationList 
                                messageBar={this.MessageBar} 
                                conversations = {this.state.conversations}
                                onConversationChange =  {this.handleConversationChange}
                                onLeaveConversation =  {this.handleLeaveConversation}
                                current_conversation_id = {this.state.current_conversation_id}
                            />
                        </Paper>
                    </Grid>
                    <Grid item xs={7} >
                        <Paper style={styles.messageListContainer} >
                            <div style={styles.messageList}>
                                <MessageList 
                                    messages = {this.state.messages} 
                                    current_conversation_type = {this.state.current_conversation_type}
                                />
                            </div>
                            <Divider light />
                            <div>
                                <MessageInput 
                                    current_conversation_id = {this.state.current_conversation_id}
                                    onNewMessage = {this.handleNewMessage}
                                />
                            </div>
                        </Paper>
                    </Grid>
                    <Grid item xs={1}>
                    </Grid>
                </Grid>

                <NewConversation 
                    messageBar={this.MessageBar} 
                    new_message_box_open={this.state.new_message_box_open}
                    onNewMessageBoxClose={this.handleNewMessageBoxClose}
                    onNewConversation={this.handleNewConversation}
                />
        </div>
    );
  }
}

export default withCookies(MessagePage);
