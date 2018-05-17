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


var styles = {
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
    },
    messageList: {
        height: "50%",
        overflow: "auto"
    },
    messageInput: {
        height: "50%"
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
            current_conversation_state: '',
            messages: '',
            suggestions: [],
            has_media: false,
            notificationsOfConversations: '',
            disableGetNotifsofConversations: true
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
        this.getConversationList((conversations) => {
            console.log(conversations)
            if (conversations.length >= 1) {
                var alreadySetNotifNum = false
                if (this.state.current_conversation_id != '') {
                    for (var i in conversations) {
                        if (conversations[i].conversation.id == this.state.current_conversation_id){
                            this.setState({
                                'current_conversation_state': conversations[i].conversation_state
                            })
                            if (!conversations[i].is_seen) {
                                this.setConversationSeen(this.state.current_conversation_id)
                                alreadySetNotifNum = true
                            }
                            break
                        }
                    }
                    this.getMessageList(this.state.current_conversation_id)  
                }
                if (!alreadySetNotifNum) {
                    this.calculateAndSetUnseenConversationNum(conversations)
                }
            } else {
                this.setCurrentConversation('', '', '')
                this.setMessageState('')
            }
            
        })         
    }

    setMessageState = (messages) => {
        this.setState({
            'messages': messages,
        });
    }

    getConversationContactList = () => {
        axios.get(
            'http://localhost:3000/conversations/getContactList',
            {
                params: {
                    'user_id': this.state.user_id
                }
            }
        ).then(response => {
            if(!response.data.result.success){
            }else{
                this.setSuggestions(response.data.result.contactList)
            }
        })
    }

    setSuggestions = (suggestions) => {
        this.setState({
            suggestions: suggestions
        })
    }

    initialize = () => {
        this.getConversationList((conversations) => {
            if (conversations.length >= 1) {
                var alreadySetNotifNum = false
                const conversation_id = conversations[0].conversation.id 
                const conversation_type = conversations[0].conversation_type
                const conversation_state = conversations[0].conversation_state

                this.setCurrentConversation(conversation_id, conversation_type, conversation_state)
                this.getMessageList(conversation_id)
                if (!conversations[0].is_seen) {
                    this.setConversationSeen(conversation_id)
                    alreadySetNotifNum = true
                }
                if (!alreadySetNotifNum) {
                    this.calculateAndSetUnseenConversationNum(conversations)
                }
            } else {
                this.setCurrentConversation('', '', '')
                this.setMessageState('')
                this.calculateAndSetUnseenConversationNum(conversations)
            }
        })    
    }

    calculateAndSetUnseenConversationNum = (conversations) => {
        var unseenConversations = 0
        for (var i in conversations) {
            if (!conversations[i].is_seen) {
                unseenConversations++
            }
        }

        if (unseenConversations == 0) {
            this.setNotificationsOfConversations('')
        } else {
            this.setNotificationsOfConversations(unseenConversations)
        }
    }

    getConversationList = (cb) => {
        axios.get(
            'http://localhost:3000/conversations',
            {
                params: {
                    'user_id': this.state.user_id
                }
            }
        ).then(response => {
            if(!response.data.result.success) {
                this.MessageBar.showSnackbar(response.data.result.error)
            }else{
                const conversations = response.data.result.conversations

                this.setConversations(conversations)

                if (cb) cb(conversations)
            }
        })
    }
    
    setConversations = (conversations) => {
        this.setState({
            'conversations': conversations,
        });
    }

    setCurrentConversation = (current_conversation_id, current_conversation_type, current_conversation_state) => {
        this.setState({
            'current_conversation_id': current_conversation_id,
            'current_conversation_type': current_conversation_type,
            'current_conversation_state': current_conversation_state
        })
    }

    getMessageList = (current_conversation_id) => {
        axios.get(
            'http://localhost:3000/messages',
            {
                params: {
                    "user_id": this.state.user_id,
                    "conversation_id": current_conversation_id
                }
            }
        ).then(response => {
            if(!response.data.result.success) {
                if (this.MessageBar != undefined) {
                    this.MessageBar.showSnackbar(response.data.result.error)
                }
            }else{
                if (response.data.result.messages == undefined) {
                    this.setMessageState('')
                } else {
                    this.setMessageState(response.data.result.messages)
                }
    
            }
        })
    }

    setNotificationsOfConversations = (notificationsOfConversations) => {
        this.setState({
            notificationsOfConversations: notificationsOfConversations
        })
    }

    setConversationSeen = (conversation_id) => {
        axios.get(
            'http://localhost:3000/conversations/setConversationSeen',
            {
                params: {
                    "user_id": this.state.user_id,
                    "conversation_id": conversation_id
                }
            }
        ).then(response => {
            if(!response.data.result.success) {
                if (this.MessageBar != undefined) {
                    this.MessageBar.showSnackbar(response.data.result.error)
                }
            } else {
                this.getConversationList((conversations) => {
                    this.calculateAndSetUnseenConversationNum(conversations)
                })
            }
        })
    }

    setConversationUnseen = (conversation_id) => {
        axios.get(
            'http://localhost:3000/conversations/setConversationUnseen',
            {
                params: {
                    "user_id": this.state.user_id,
                    "conversation_id": conversation_id
                }
            }
        ).then(response => {
            if(!response.data.result.success) {
                if (this.MessageBar != undefined) {
                    this.MessageBar.showSnackbar(response.data.result.error)
                }
            }
        })
    }

    handleNewMessageBoxOpen = e => {
        this.getConversationContactList()
        this.setState({new_message_box_open: true });
    }

    handleNewMessageBoxClose = e => {
        this.setState({new_message_box_open: false });
    }

    handleNewConversation = (new_conversation_ID, new_conversation_type, new_conversation_state) => {
        this.getConversationList()
        this.setCurrentConversation(new_conversation_ID, new_conversation_type, new_conversation_state)
        this.getMessageList(new_conversation_ID)
    }

    setSawMessagesUntil = (conversation_id, message_id) => {
        axios.get(
            'http://localhost:3000/messages/setSawMessagesUntil',
            {
                params: {
                    "user_id": this.state.user_id,
                    "conversation_id": conversation_id,
                    "message_id": message_id
                }
            }
        ).then(response => {
            if(!response.data.result.success) {
                if (this.MessageBar != undefined) {
                    this.MessageBar.showSnackbar(response.data.result.error)
                }
            }
        })
    }

    handleConversationChange = (current_conversation_id) => {
        if (current_conversation_id == this.state.current_conversation_id) {
            return
        }
        var conversations = this.state.conversations
        var current_conversation_type 
        var current_conversation_state
        var is_seen
        for (var i in conversations) {
            var conversation = conversations[i].conversation
            if (conversation.id == current_conversation_id) {
                current_conversation_type = conversation.conversation_type
                current_conversation_state = conversation.conversation_state
                is_seen = conversations[i].is_seen
                break
            }
        }

        if (!is_seen) {
            this.setConversationSeen(current_conversation_id)
        }
        
        const messages = this.state.messages
        if (messages.length > 0) {
            this.setSawMessagesUntil(this.state.current_conversation_id, messages[messages.length-1].id)
            this.changeSawMessagesUntilInState(messages[messages.length-1])
        }
        
        this.setCurrentConversation(current_conversation_id, current_conversation_type, current_conversation_state)
        this.getMessageList(current_conversation_id)
    }

    handleLeaveConversation = () => {
        this.initialize()
    }

    changeSawMessagesUntilInState = (newMessage) => {
        var conversations = this.state.conversations
        for (var i in conversations) {
            if (conversations[i].conversation.id == newMessage.conversation_id) {
                for (var j in conversations[i].conversation_participants) {
                    if (conversations[i].conversation_participants[j].id == this.state.user_id) {
                        conversations[i].conversation_participants[j].saw_messages_until = newMessage.created_at
                        this.setConversations(conversations)
                    }
                }
            }
        }
    }

    handleNewMessage = (newMessage) => {
        this.changeSawMessagesUntilInState(newMessage)
        
        this.getMessageList(this.state.current_conversation_id)
        this.setConversationUnseen(this.state.current_conversation_id)
    }

    setHasMediaState = (has_media) => {
        this.setState({has_media: has_media})
    }
    
    changeMessageInputLayout = () => {
        if (this.state.has_media) {
            styles.messageInput = {
                height: '50%'
            }
            styles.messageList = {
                height: "50%",
                overflow: "auto"
            }
        } else {
            styles.messageInput = {
                height: "13%"
            }
            styles.messageList = {
                height: "87%",
                overflow: "auto"
            }
        }
    }

    getSawMessagesUntil = (conversation_id) => {
        console.log(conversation_id)
        const conversations = this.state.conversations
        for (var i in conversations) {
            const conversation = conversations[i]
            if (conversation.conversation.id == conversation_id) {
                const conversation_participants = conversation.conversation_participants
                for (var j in conversation_participants) {
                    if (conversation_participants[j].id == this.state.user_id) {
                        return conversation_participants[j].saw_messages_until
                    }
                }
            }
        }
    }

    render () {
        this.changeMessageInputLayout()

        return (
            <div>
                <MessageBar ref={instance => { this.MessageBar = instance; }}/>

                <NavBar 
                    notificationsOfConversations={this.state.notificationsOfConversations}
                    disableGetNotifsofConversations={this.state.disableGetNotifsofConversations}
                />

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
                            <Grid style={styles.messageList}>
                                <MessageList 
                                    messages = {this.state.messages}
                                    saw_messages_until = {this.getSawMessagesUntil(this.state.current_conversation_id)}
                                    current_conversation_type = {this.state.current_conversation_type}
                                />
                            </Grid>
                            <Grid>
                                <Divider light />
                            </Grid>
                            <Grid style={styles.messageInput}>
                                <MessageInput 
                                    messageBar = {this.MessageBar}
                                    current_conversation_id = {this.state.current_conversation_id}
                                    current_conversation_state = {this.state.current_conversation_state}
                                    onNewMessage = {this.handleNewMessage}
                                    messageBar={this.MessageBar} 
                                    setHasMediaState={this.setHasMediaState}
                                />
                            </Grid>
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
                    suggestions={this.state.suggestions}
                />
        </div>
    );
  }
}

export default withCookies(MessagePage);
