import React, {Component, Fragment} from "react";
import Divider from 'material-ui/Divider';
import Conversation from './Conversation'
import List from 'material-ui/List';

class ConversationList extends Component {

    constructor(props) {

        super(props);

        this.state = {
        }
        
    }

    render () {
        const conversations = this.props.conversations
        const conversationList = conversations.map((conversation, index) =>
            <div key={conversation.conversation.id}>
                <Conversation conversation = {conversation}/>
                <li>
                    <Divider inset />
                </li>
            </div>
        );


        return (
            <List>
                {conversationList}
            </List>
        )
    }

}

export default ConversationList