import React, {Component, Fragment} from "react";
import List from 'material-ui/List';
import Message from './Message'
import Typography from 'material-ui/Typography';
import { color } from "material-ui/colors";

const styles = {
    unreadLine: {
        width           : '57%',
        marginLeft      : "auto",
        marginRight     : "auto",
        marginBottom    : 40,
        marginTop       : 40,
        color           : '#9E9E9E',
        padding         : 10
    }
}

class MessageList extends Component {

    constructor(props) {
        super(props);
    }

    showUnreadLine = () => {
        return (
            <Typography style={styles.unreadLine}>
                ------------------------------------------ Unread Messages ------------------------------------------
            </Typography>
        )
    }

    render () {
        const messages = this.props.messages
        var messageList = ''
        if (Array.isArray(messages)) {
            var addUnreadLine = false
            messageList = messages.map((message) => {
                if (this.props.saw_messages_until < message.created_at && !addUnreadLine) {
                    console.log("this " + this.props.saw_messages_until)
                    console.log("that " + message.created_at)
                    addUnreadLine = true
                    return (
                        <Fragment>
                            {this.showUnreadLine()}
                            <Message key={message.id}
                                message = {message}
                                current_conversation_type = {this.props.current_conversation_type}
                            />
                        </Fragment>
                    )
                } else {
                    return (
                        <Message key={message.id}
                            message = {message}
                            current_conversation_type = {this.props.current_conversation_type}
                        />
                    )
                }
            })
        }
        return (
            <Fragment>
                {messageList}
            </Fragment>
        )   
    }
}

export default MessageList