import React, {Component, Fragment} from "react";
import MoreVertIcon from '@material-ui/icons/MoreVert';
import IconButton from 'material-ui/IconButton';
import Menu, { MenuItem } from 'material-ui/Menu';
import axios from 'axios';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';

const options = [
    'Leave Conversation'
];

class ConversationActions extends Component {
    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    }

    constructor(props) {

        super(props);

        const { cookies } = this.props;

        this.state = {
            anchorEl: null,
            user_id: cookies.get('user_id'),
        }
        
    }

    handleClick = e => {
        this.setState({ anchorEl: e.currentTarget });
    };
    
    handleClose = () => {
        this.setState({ anchorEl: null });
    };

    handleOption = (e, option) => {
        if (option == options[0]) {
            this.handleLeaveConversation()
        }

        this.setState({ anchorEl: null });
    }

    handleLeaveConversation = () => {
        axios.delete(
            'http://localhost:3000/messages/',
            {
                params: {
                    "user_id": this.state.user_id,
                    "conversation_id": this.props.current_conversation_id,
                }
            }
        ).then(response => {
            if(!response.data.result.success){
            }else{
                console.log('ooooooooooooooo')
                // this.props.onNewMessage()
            }
        })
    }

    render() {

        return (
            <Fragment>
                <IconButton
                    aria-label="More"
                    aria-owns={this.state.anchorEl ? 'long-menu' : null}
                    aria-haspopup="true"
                    onClick={this.handleClick}
                >
                <MoreVertIcon />
                </IconButton>
                <Menu
                    id="long-menu"
                    anchorEl={this.state.anchorEl}
                    open={Boolean(this.state.anchorEl)}
                    onClose={this.handleClose}
                    PaperProps={{
                        style: {
                        maxHeight: options.length * 30 * 4.5,
                        width: 200,
                        },
                    }}
                >
                    {options.map(option => (
                        <MenuItem key={option} onClick={(e) => this.handleOption(e, option)}>
                            {option}
                        </MenuItem>
                    ))}
                </Menu>
            </Fragment>
        )
    }
}

export default withCookies(ConversationActions);