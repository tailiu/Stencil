import React, {Component} from "react";
import Dialog, {
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
} from 'material-ui/Dialog';
import Button from 'material-ui/Button';
import axios from 'axios';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
import NewConversationSearchUser from './NewConversationSearchUser'

const styles = {
    height: 300
}

class NewConversation extends Component {
    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    };

    constructor(props) {
        super(props);

        const { cookies } = this.props;

        this.state = {
            user_handle: cookies.get('user_handle'),
            selectedItem: []
        }
    }

    clearMessageTo = () => {
        this.setState({
            selectedItem: []
        })
    }

    setSelectedItem = (selectedItem) => {
        this.setState({
            selectedItem: selectedItem
        })
    }

    handleNewConversation = e => {
        var participants = {
            participants: [],
            conversation_creator: ''
        }
        var selectedItem = this.state.selectedItem

        for (var i in selectedItem) {
            const handle = selectedItem[i].split('@')[1]
            participants.participants.push(handle)
        }

        participants.participants.push(this.state.user_handle)
        participants.conversation_creator = this.state.user_handle

        axios.get(
            'http://localhost:3000/conversations/new',
            {
                params: {
                    'participants': participants
                }
            }
        ).then(response => {
            if(!response.data.result.success){
                this.props.messageBar.showSnackbar(response.data.result.error)
            }else{
                if (response.data.result.message != "") {
                    this.props.messageBar.showSnackbar(response.data.result.message)
                }
                
                this.clearMessageTo()

                const conversation = response.data.result.conversation
                const conversation_state = response.data.result.conversation_state

                this.handleNewMessageBoxClose()
                this.props.onNewConversation(conversation.id, conversation.conversation_type, conversation_state )
            }
        })

    }

    handleNewMessageBoxClose = () => {
        this.clearMessageTo()
        this.props.onNewMessageBoxClose()
    }

    render () {
        const open = this.props.new_message_box_open

        return (
            <Dialog
                open={open}
                onClose={this.handleNewMessageBoxClose}
                aria-labelledby="form-dialog-title"
                fullWidth
            >
                <DialogTitle id="form-dialog-title">New Message</DialogTitle>
                <DialogContent style={styles}>
                    <NewConversationSearchUser 
                        suggestions={this.props.suggestions}
                        selectedItem={this.state.selectedItem}
                        setSelectedItem={this.setSelectedItem}
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={this.handleNewMessageBoxClose} color="primary">
                        Cancel
                    </Button>
                    <Button onClick={this.handleNewConversation} color="primary">
                        New Message
                    </Button>
                </DialogActions>
            </Dialog>
        )
    }
}

export default withCookies(NewConversation);




