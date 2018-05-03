import React, {Component, Fragment} from "react";
import Avatar from 'material-ui/Avatar';
import List, { ListItem, ListItemText, } from 'material-ui/List';
import Moment from 'moment';
import axios from 'axios';

class Conversation extends Component {
    constructor(props) {

        super(props);

        this.state = {
        }
        
    }

    handleClick = e => {
        this.props.onConversationChange(this.props.conversation.conversation.id)
    } 

    getTitleForConversation = () => {
        const conversation_participants = this.props.conversation.conversation_participants

        var conversationTitle = ''
        for (var i in conversation_participants) {
            conversationTitle += '@' + conversation_participants[i].handle + ' '
        }

        return conversationTitle
    }

    getLatestUpdatedDateForConversation = () => {
        return Moment(this.props.conversation.conversation.updated_at).format('MMMM Do, YYYY - h:mm A')
    }

    render() {
        const title = this.getTitleForConversation()
        const latestUpdatedDate = this.getLatestUpdatedDateForConversation()

        return (
            <ListItem onClick={this.handleClick}>        
                <Avatar src={require('../Assets/Images/user_icon.png')} />
                <ListItemText primary={title} secondary={latestUpdatedDate} />
            </ListItem>
        )
    }
}

export default Conversation