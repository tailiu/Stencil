import React, {Component, Fragment} from "react";
import Divider from 'material-ui/Divider';
import Avatar from 'material-ui/Avatar';
import List, { ListItem, ListItemText, } from 'material-ui/List';


class ConversationList extends Component {

    constructor(props) {

        super(props);

        this.state = {
        }
        
    }

    handleClick = e => {
        console.log(e.target.value.primary)
    };

    formConversationList = conversations => {
        var conversationList = [];

        for (var i in conversations) {
            var conversationTitle = ''
            for (var j in conversations[i].conversation_participants) {
                conversationTitle += '@' + conversations[i].conversation_participants[j].handle + ' '
            }
            conversationList.push(conversationTitle);
        }

        const list = conversationList.map((title, index) =>
            <div>
                <ListItem onClick={this.handleClick}>
                    <Avatar src={require('../Assets/Images/user_icon.png')} />
                    <ListItemText primary={title} secondary="Jan 9, 2014" />
                </ListItem>
                <li>
                    <Divider inset />
                </li>
            </div>
        );

        return (
            <div>
                {list}
            </div>
        );

    }

    render () {
        const conversationList = this.formConversationList(this.props.conversations)
        
        return (
            <List>
                {conversationList}
            </List>
        )
    }

}

export default ConversationList