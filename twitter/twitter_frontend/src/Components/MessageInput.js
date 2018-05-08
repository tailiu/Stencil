import React, {Component} from "react";
import axios from 'axios';
import TextField from 'material-ui/TextField';
import Button from 'material-ui/Button';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';

var styles = {
    inputContainer: {
        marginTop: 25
    },
    messagesInput: {
        width: "80%",
        float: "left",
        marginLeft: 30,
        backgroundColor: "#fff",
    }
}

class MessageInput extends Component {

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    }

    constructor(props) {
        super(props);

        const { cookies } = this.props;

        this.state = {
            value: '',
            user_id: cookies.get('user_id')
        }
    }

    handleChange = (e) => {
       this.setState({
           value: e.target.value
       })
    }

    handleNewMessage = () => {
        axios.get(
            'http://localhost:3000/messages/new',
            {
                params: {
                    "user_id": this.state.user_id,
                    "conversation_id": this.props.current_conversation_id,
                    "content": this.state.value
                }
            }
        ).then(response => {
            if(!response.data.result.success){
            }else{
                this.setState({value : ''});
                this.props.onNewMessage()
            }
        })
    }

    catchReturn = (e) => {
        if (e.key === 'Enter' && this.state.value != '' && this.props.current_conversation_id != '') {
            this.handleNewMessage()
        }
        
    }

    render() {
        var disabled = true

        if (this.state.value != '' && this.props.current_conversation_id != '') {
            disabled = false
        }
        
        styles.sendMessageButton = {
            backgroundColor: disabled ? '#BBDEFB' : "#00aced",
            color: "#fff",
            variant: "raised",
            display: "inline-block",
            float: "right",
            marginRight: 30
        }

        return (
            <div style={styles.inputContainer}>
                <TextField
                    id="message"
                    margin="normal"
                    fullWidth
                    style={styles.messagesInput}
                    value={this.state.value}
                    onChange={this.handleChange}
                    onKeyPress={this.catchReturn}
                />
                <Button style={styles.sendMessageButton} onClick={this.handleNewMessage} color="primary" disabled={disabled}>
                    Send
                </Button>
            </div>
        )
    }
}

export default withCookies(MessageInput);