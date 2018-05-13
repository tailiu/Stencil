import React, {Component} from "react";
import MessageInputAllow from './MessageInputAllow'
import MessageInputBlock from './MessageInputBlock'

class MessageInput extends Component {

    constructor(props) {
        super(props)
    }
 
    render() {
        if (this.props.current_conversation_state == 'blocked') {
            return (
                <MessageInputBlock 
                />
            )
        } else {
            return (
                <MessageInputAllow 
                    current_conversation_id = {this.props.current_conversation_id}
                    onNewMessage = {this.props.onNewMessage}
                    messageBar={this.props.messageBar} 
                />
            )
        }
    }
}

export default MessageInput;