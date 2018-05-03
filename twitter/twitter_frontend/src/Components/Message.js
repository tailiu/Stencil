import Grid from 'material-ui/Grid';
import React, {Component, Fragment} from "react";
import List, { ListItem, ListItemText, } from 'material-ui/List';
import Moment from 'moment';

class Message extends Component {

    constructor(props) {
        super(props);
    }

    getLatestUpdatedDateForMessage = () => {
        return Moment(this.props.message.updated_at).format('MMMM Do, YYYY - h:mm A');
    }

    render () {
        const content = this.props.message.content;
        const updatedDate = this.getLatestUpdatedDateForMessage();

        return (
            <div>
                <ListItem>
                    <ListItemText
                        primary={content}
                        secondary={updatedDate}
                    />
                </ListItem>
            </div>
        )   
    }
}

export default Message