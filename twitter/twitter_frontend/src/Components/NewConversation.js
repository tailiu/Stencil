import React, {Component} from "react";
import Dialog, {
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
} from 'material-ui/Dialog';
import TextField from 'material-ui/TextField';
import Button from 'material-ui/Button';
import axios from 'axios';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';

class NewConversation extends Component {
    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    };

    constructor(props) {
        super(props);

        const { cookies } = this.props;

        this.state = {
            message_to: '',
            user_handle: cookies.get('user_handle')
        }
    }

    updateMessageTo = e => {
        this.setState({
            message_to: e.target.value
        })
    }

    validateInput = () => {
        if (this.state.message_to.indexOf("@") == -1) {
            return false
        }
        return true
    }

    handleNewConversation = e => {
        if(!this.validateInput()){
            this.props.messageBar.showSnackbar("Please input valid user handles (Ex: @userhandle)")
            return
        }

        const raw_data = this.state.message_to.split('@')
        raw_data.shift()

        const participants = []
        for (var i in raw_data) {
            raw_data[i] = raw_data[i].replace(/\s/g,''); // replace all spaces in handles
            if (raw_data[i] != this.state.user_handle) {
                participants.push(raw_data[i])
            }
        }
        participants.push(this.state.user_handle)

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
                this.setState({
                    message_to: ''
                })

                const conversation = response.data.result.conversation
                this.props.onNewMessageBoxClose()
                this.props.onNewConversation(conversation.id, conversation.conversation_type)
            }
        })

    }

    render () {
        const open = this.props.new_message_box_open

        return (
            <Dialog
                open={open}
                onClose={this.props.onNewMessageBoxClose}
                aria-labelledby="form-dialog-title"
                fullWidth
                >
                <DialogTitle id="form-dialog-title">New Message</DialogTitle>
                <DialogContent>
                    <TextField
                        autoFocus
                        margin="dense"
                        id="tweet"
                        label="Send message to"
                        type="email"
                        placeholder="Example: @tai @zain"
                        value={this.state.message_to}
                        onChange={this.updateMessageTo}
                        fullWidth
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={this.props.onNewMessageBoxClose} color="primary">
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