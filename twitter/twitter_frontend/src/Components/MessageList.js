import Grid from 'material-ui/Grid';
import React, {Component} from "react";
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
                <div key={message.id}>
                    <Message message = {message} />
                </div>
            );
        } 

        return (
            <List dense={true}>
                {messageList}
            </List>
        )   
    }
}

export default MessageList