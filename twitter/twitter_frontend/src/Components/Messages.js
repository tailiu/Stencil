import React, {Component, Fragment} from "react";
import TextField from 'material-ui/TextField';
import Grid from 'material-ui/Grid';
import Divider from 'material-ui/Divider';
import Card, { CardContent, CardHeader } from 'material-ui/Card';
import NavBar from './NavBar';
import Avatar from 'material-ui/Avatar';
import List, { ListItem, ListItemText, } from 'material-ui/List';

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
    }
};

function generate(element) {
    return [0, 1, 2].map(value =>
      React.cloneElement(element, {
        key: value,
      }),
    );
  }

class Messages extends Component {

    constructor(props) {

        super(props);

        this.state = {
        }
    }

    handleSubmit = e => {
        e.preventDefault();
    }

  render () {
    return (
        <Fragment>
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
                            />
                            <hr />
                            <CardContent>
                                <Grid container direction="row" spacing={8} align="left">
                                    <Grid item xs={4}>
                                        <List>
                                            <ListItem>
                                            <Avatar
                                            src={require('../Assets/Images/user_icon.png')}
                                            />
                                            <ListItemText primary="Tai Cow" secondary="Jan 9, 2014" />
                                            </ListItem>
                                            <li>
                                            <Divider inset />
                                            </li>
                                            <ListItem>
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
                                            </li>
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

export default Messages;
