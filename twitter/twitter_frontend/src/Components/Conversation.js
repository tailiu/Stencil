import React, {Component} from "react";
import Avatar from 'material-ui/Avatar';
import { 
    ListItemText, 
    ListItemSecondaryAction,
} from 'material-ui/List';
import Moment from 'moment';
import axios from 'axios';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import ConversationActions from './ConversationActions'
import { MenuItem } from 'material-ui/Menu';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';

const styles = {
    menuItem: {
        padding: 20
    }
}

class Conversation extends Component {

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    }

    constructor(props) {

        super(props);

        const { cookies } = this.props;

        this.state = {
            user_handle: cookies.get('user_handle'),
        };
        
    }

    handleClick = e => {
        this.props.onConversationChange(
            this.props.conversation.conversation.id, 
            this.props.conversation.conversation.conversation_type,
            this.props.conversation.conversation_state
        )
    } 

    getTitleForConversation = () => {
        const conversation_participants = this.props.conversation.conversation_participants

        var conversationTitle = ''
        var handle = ''

        if (conversation_participants.length == 1) {
            conversationTitle += 'You'
        } else if (conversation_participants.length == 2) {
            for (var i in conversation_participants) {
                handle = conversation_participants[i].handle
                if (this.state.user_handle == handle) {
                    continue
                } else {
                    conversationTitle += '@' + conversation_participants[i].handle + ' '
                }
            }
        } else {
            for (var i in conversation_participants) {
                handle = conversation_participants[i].handle
                if (this.state.user_handle == handle) {
                    conversationTitle += 'You '
                } else {
                    conversationTitle += '@' + conversation_participants[i].handle + ' '
                }
                if (i == conversation_participants.length - 2) {
                    conversationTitle += 'and '
                }
            }
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
            <MenuItem   button 
                        onClick={this.handleClick} 
                        selected={this.props.selected === this.props.conversation.conversation.id}
                        style={styles.menuItem}
            >        
                <Avatar src={require('../Assets/Images/user_icon.png')} />
                <ListItemText primary={title} secondary={latestUpdatedDate} />
                <ListItemSecondaryAction>
                    <ConversationActions
                        messageBar = {this.props.messageBar} 
                        conversationID = {this.props.conversation.conversation.id}
                        onLeaveConversation = {this.props.onLeaveConversation}
                    />
                </ListItemSecondaryAction>
            </MenuItem>
        )
    }
}

export default withCookies(Conversation);