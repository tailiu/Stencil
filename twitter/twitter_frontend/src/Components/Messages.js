import React, {Component, Fragment} from "react";
import TextField from 'material-ui/TextField';
import Grid from 'material-ui/Grid';
import Divider from 'material-ui/Divider';
import Card, { CardContent, CardHeader } from 'material-ui/Card';
import NavBar from './NavBar';
import Avatar from 'material-ui/Avatar';
import List, { ListItem, ListItemText, } from 'material-ui/List';
import Button from 'material-ui/Button';
import Dialog, {
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
} from 'material-ui/Dialog';
import MessageBar from './MessageBar';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
import axios from 'axios';

const styles = {
    grid : {
        container : {
            marginTop: 80
        }
    },
    messages: {
        input: {
            marginTop: 20,
            marginLeft: 20,
            width: "90%"
        }
    },
    tweetButton: {
        backgroundColor: "#00aced",
        color: "#fff"
    }
};

function generate(element) {
    return [0, 1, 2].map(value =>
      React.cloneElement(element, {
        key: value,
      }),
    );
}

class ConversationOverview extends Component {

    constructor(props) {

        super(props);

        this.state = {
        }
    }

    render () {
        return (
            <div>
                <ListItem>
                    <Avatar src={require('../Assets/Images/user_icon.png')} />
                <ListItemText primary="Tai Cow" secondary="Jan 9, 2014" />
                </ListItem>
                <li>
                    <Divider inset />
                </li>
            </div>
        )
    }

}


class Messages extends Component {

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    };

    constructor(props) {

        super(props);

        const { cookies } = this.props;

        this.state = {
            user_id : cookies.get('user_id'),
            new_message_box_open: false,
            message_to: ''
        }
    }

    handleNewMessageBoxOpen = e => {
        this.setState({new_message_box_open: true });
    }

    handleNewMessageBoxClose = e => {
        this.setState({new_message_box_open: false });
    }

    updateMessageTo = e => {
        this.setState({
            message_to: e.target.value
        })
    }

    validateInput = e => {
        return true
    }

    handleNewMessage = e => {
        if(!this.validateInput()){
            this.MessageBar.showSnackbar("Please input valid user handles")
            return
        }

        const raw_data = this.state.message_to.split('@')
        raw_data.shift()

        const participants = []
        for (var i in raw_data) {
            raw_data[i] = raw_data[i].replace(/\s/g,''); // replace all spaces in handles
            participants.push(raw_data[i])
        }

        axios.get(
            'http://localhost:3000/conversations/new',
            {
                params: {
                    'id': this.state.user_id,
                    'participants': participants
                }
            }
        ).then(response => {
            // console.log(response)
            if(!response.data.result.success){
                this.MessageBar.showSnackbar(response.data.result.error.message)
            }else{

            }
        })

    }

    render () {
        return (
            <Fragment>
                <MessageBar ref={instance => { this.MessageBar = instance; }}/>
                <NavBar />
                <Grid style={styles.grid.container} container spacing={24} >
                    
                    <Grid item xs={1}>
                    </Grid>
                    <Grid item xs={10}>
                        <Grid container direction="column" align="left">
                            <Grid item>
                            <Card>
                                <CardHeader
                                    title="Messages"

                                    action={
                                        <Button 
                                        style={styles.tweetButton} 
                                        onClick={this.handleNewMessageBoxOpen}>New Message</Button>
                                    }
                                />

                                <Dialog
                                    open={this.state.new_message_box_open}
                                    onClose={this.handleNewMessageBoxClose}
                                    aria-labelledby="form-dialog-title"
                                    >
                                    <DialogTitle id="form-dialog-title">New Message</DialogTitle>
                                    <DialogContent>
                                        <DialogContentText>
                                        {/* What's on your mind? */}
                                        </DialogContentText>
                                        <TextField
                                        autoFocus
                                        margin="dense"
                                        id="tweet"
                                        label="Send message to"
                                        type="email"
                                        value={this.state.message_to}
                                        onChange={this.updateMessageTo}
                                        fullWidth
                                        />
                                    </DialogContent>
                                    <DialogActions>
                                        <Button onClick={this.handleNewMessageBoxClose} color="primary">
                                            Cancel
                                        </Button>
                                        <Button onClick={this.handleNewMessage} color="primary">
                                            New Message
                                        </Button>
                                    </DialogActions>
                                </Dialog>


                                <hr />

                                <CardContent>
                                    <Grid container direction="row" spacing={8} align="left">
                                        <Grid item xs={4}>
                                            <List>
                                                <ConversationOverview />

                                                {/* <ListItem>
                                                <Avatar
                                                src={require('../Assets/Images/user_icon.png')}
                                                />
                                                <ListItemText primary="Miro Pasta" secondary="Jan 9, 2014" />
                                                </ListItem>
                                                <li>
                                                <Divider inset />
                                                </li>

                                                <ListItem>
                                                <Avatar
                                                src={require('../Assets/Images/user_icon.png')}
                                                />
                                                <ListItemText primary="Major Tom" secondary="Jan 9, 2014" />
                                                </ListItem>
                                                <li>
                                                <Divider inset />
                                                </li> */}
                                            </List>
                                        </Grid>
                                        <Grid item xs={8} >
                                            <Grid container direction="column">
                                                <Grid item>
                                                    <List dense={true}>
                                                        {generate(
                                                        <ListItem>
                                                            <ListItemText
                                                            primary="Miro: Hey!"
                                                            secondary="Jan 9, 2017"
                                                            />
                                                        </ListItem>,
                                                        )}
                                                    </List>
                                                </Grid>
                                            </Grid>
                                            <Grid>
                                                <TextField
                                                    id="message"
                                                    label="Message"
                                                    margin="normal"
                                                    fullWidth
                                                    style={styles.messages.input}
                                                    // onChange={this.handleChange}
                                                />
                                            </Grid>
                                        </Grid>
                                    </Grid>                                
                                </CardContent>
                                </Card>
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item xs={1}>
                    </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default withCookies(Messages);
