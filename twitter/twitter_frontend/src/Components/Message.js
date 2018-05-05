import React, {Component, Fragment} from "react";
import { ListItem, ListItemText, } from 'material-ui/List';
import Moment from 'moment';
import Avatar from 'material-ui/Avatar';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';

class Message extends Component {

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    }

    constructor(props) {
        super(props);

        const { cookies } = this.props;

        this.state = {
            user_id: cookies.get('user_id'),
        };
    }

    getLatestUpdatedDateForMessage = () => {
        return Moment(this.props.message.updated_at).format('MMMM Do, YYYY - h:mm A');
    }

    setStyle = () => {
        const message = this.props.message

        var styles = {
            listContainer: {
                float       : 'none', 
                width       : '20vw',
                marginLeft  : 0,
                marginRight : 0
            },
            listItem: {
                whiteSpace: 'normal',
                wordWrap: 'break-word'
            }
        }

        if (message.user_id == this.state.user_id) {
            styles.listContainer.marginLeft = 'auto'
        }

        return styles
    }

    render () {
        const content = this.props.message.content;
        const updatedDate = this.getLatestUpdatedDateForMessage();

        const styles = this.setStyle()

        return (
            <div>
                <ListItem style={styles.listContainer}>
                    <Avatar src={require('../Assets/Images/user_icon.png')} />
                    <ListItemText style={styles.listItem}
                        primary={content}
                        secondary={updatedDate}
                    />
                </ListItem>
            </div>
        )   
    }
}

export default withCookies(Message);