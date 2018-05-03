import Grid from 'material-ui/Grid';
import React, {Component, Fragment} from "react";
import axios from 'axios';
import TextField from 'material-ui/TextField';
import Button from 'material-ui/Button';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';

const styles = {
    messages: {
        input: {
            marginTop: 20,
            marginLeft: 20,
            width: "90%"
        }
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
            user_id: cookies.get('user_id'),
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
        if (e.key === 'Enter') {
            this.handleNewMessage()
        }
        
    }

    render() {
        return (
            <div>
                <Grid>
                    <TextField
                        id="message"
                        margin="normal"
                        fullWidth
                        style={styles.messages.input}
                        value={this.state.value}
                        onChange={this.handleChange}
                        onKeyPress={this.catchReturn}
                    />
                    <Button onClick={this.handleNewMessage} color="primary">
                        Send
                    </Button>
                </Grid>
            </div>
        )
    }
}

export default withCookies(MessageInput);