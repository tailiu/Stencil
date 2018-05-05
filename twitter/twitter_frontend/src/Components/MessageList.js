import Grid from 'material-ui/Grid';
import React, {Component} from "react";
import List, { ListItem, ListItemText, } from 'material-ui/List';
import Message from './Message'


class MessageList extends Component {

    constructor(props) {
        super(props);
    }

    render () {
        const messages = this.props.messages
        var messageList = <div></div>
        if(Array.isArray(messages)) {
            messageList = messages.map((message) =>
                <div key={message.id}>
                    <Message message = {message} />
                </div>
            );
        } 

        return (
            <div>
                <Grid container direction="column">
                    <Grid item>
                        <List dense={true}>
                            {messageList}
                        </List>
                    </Grid>
                </Grid>
            </div>
        )   
    }
}

export default MessageList