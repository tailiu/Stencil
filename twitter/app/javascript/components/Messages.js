import React, {Component, Fragment} from "react";

import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import MenuItem from 'material-ui/Menu/MenuItem';
import TextField from 'material-ui/TextField';

import Paper from 'material-ui/Paper';
import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';
import Divider from 'material-ui/Divider';

import Button from 'material-ui/Button';
import Card, { CardActions, CardContent, CardHeader } from 'material-ui/Card';

import NavBar from './NavBar';

import Avatar from 'material-ui/Avatar';

import Collapse from 'material-ui/transitions/Collapse';
import IconButton from 'material-ui/IconButton';
import red from 'material-ui/colors/red';

import UserIcon from 'images/user_icon.png';

import List, {
    ListItem,
    ListItemAvatar,
    ListItemIcon,
    ListItemSecondaryAction,
    ListItemText,
  } from 'material-ui/List';

const styles = {
    grid : {
        container : {
            marginTop: 80
        }
    },
    messages: {
        input: {
            marginLeft: 22,
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

class Home extends Component {

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
                                            src={UserIcon}
                                            />
                                            <ListItemText primary="Tai Cow" secondary="Jan 9, 2014" />
                                            </ListItem>
                                            <li>
                                            <Divider inset />
                                            </li>
                                            <ListItem>
                                            <Avatar
                                            src={UserIcon}
                                            />
                                            <ListItemText primary="Miro Pasta" secondary="Jan 9, 2014" />
                                            </ListItem>
                                            <li>
                                            <Divider inset />
                                            </li>
                                            <ListItem>
                                            <Avatar
                                            src={UserIcon}
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

export default withStyles(styles)(Home);