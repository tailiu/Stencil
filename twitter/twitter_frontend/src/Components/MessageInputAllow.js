import React, {Component} from "react";
import axios from 'axios';
import TextField from 'material-ui/TextField';
import Button from 'material-ui/Button';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
import Grid from 'material-ui/Grid';
import Card, { CardActions } from 'material-ui/Card';
import FileUpload from '@material-ui/icons/FileUpload';

var styles = {
    inputContainer: {
        height: '100%'
    },
    messagesInput: {
        width: "75%",
        marginLeft: 5,
        marginRight: 10
    }
}

class MessageInputAllow extends Component {

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
            if(!response.data.result.success) {
                this.props.messageBar.showSnackbar(response.data.result.error)
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
            variant: "raised"
        }

        return (
            <Card style={styles.inputContainer}>
                <CardActions>
                    <TextField
                        margin="normal"
                        style={styles.messagesInput}
                        value={this.state.value}
                        onChange={this.handleChange}
                        onKeyPress={this.catchReturn}
                    />
                    <Button size="small" variant="raised" color="default">
                        Upload
                        <FileUpload  />
                    </Button>
                    <Button size="large" style={styles.sendMessageButton} onClick={this.handleNewMessage} color="primary" disabled={disabled}>
                        Send
                    </Button>
                </CardActions>
            </Card>
        )
    }
}

export default withCookies(MessageInputAllow);