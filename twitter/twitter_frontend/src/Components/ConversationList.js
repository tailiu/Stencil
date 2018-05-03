import React, {Component, Fragment} from "react";
import Divider from 'material-ui/Divider';
import Conversation from './Conversation'
import List from 'material-ui/List';
// import Subheader from 'material-ui/Subheader';

class ConversationList extends Component {

    constructor(props) {

        super(props);

        this.state = {
        }
        
    }

    render () {
        const conversations = this.props.conversations
        const conversationList = conversations.map((conversation) =>
            <div key={conversation.conversation.id}>
                <Conversation conversation = {conversation} onConversationChange = {this.props.onConversationChange}/>
                <li>
                    <Divider inset />
                </li>
            </div>
        );


        return (
            <List>
                {/* <Subheader inset={true}>Folders</Subheader> */}
                {conversationList}
            </List>
        )
    }

}

export default ConversationList