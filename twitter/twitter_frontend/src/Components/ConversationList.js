import React, {Component} from "react";
import Divider from 'material-ui/Divider';
import Conversation from './Conversation'
import List, {ListSubheader} from 'material-ui/List';
import { MenuList, MenuItem } from 'material-ui/Menu';


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
            <MenuList subheader={<ListSubheader>Conversations</ListSubheader>} >
                {conversationList}
            </MenuList>
        )
    }

}

export default ConversationList