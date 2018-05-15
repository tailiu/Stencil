import React, {Component, Fragment} from "react";
import List from 'material-ui/List';
import Message from './Message'

class MessageList extends Component {

    constructor(props) {
        super(props);
    }

    render () {
        const messages = this.props.messages
        var messageList = ''
        if(Array.isArray(messages)) {
            messageList = messages.map((message) =>
                <Message key={message.id}
                    message = {message}
                    current_conversation_type = {this.props.current_conversation_type}
                />
            );
        } 

        return (
            <Fragment>
                {messageList}
            </Fragment>
        )   
    }
}

export default MessageList