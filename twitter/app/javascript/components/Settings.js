import React, {Component, Fragment} from "react";

import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import MenuItem from 'material-ui/Menu/MenuItem';
import TextField from 'material-ui/TextField';

import Paper from 'material-ui/Paper';
import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';

import Button from 'material-ui/Button';
import Card, { CardActions, CardContent, CardHeader } from 'material-ui/Card';

import NavBar from './NavBar';

import Avatar from 'material-ui/Avatar';

import Collapse from 'material-ui/transitions/Collapse';
import IconButton from 'material-ui/IconButton';
import red from 'material-ui/colors/red';

import SearchIcon from 'images/search_icon.png';
import TwitterLogo from 'images/Twitter_Logo_Blue.png';
import ReplyIcon from 'images/reply_icon.png';
import FavIcon from 'images/fav_icon.png';
import RetweetIcon from 'images/retweet_icon.png';
import FollowIcon from 'images/follow_icon.png';
import Checkbox from 'material-ui/Checkbox';

import { FormGroup, FormControlLabel } from 'material-ui/Form';

import List, { ListItem, ListItemIcon, ListItemText } from 'material-ui/List';
import Divider from 'material-ui/Divider';

const styles = {
    grid : {
        container : {
            marginTop: 80,
            height: 200
        }
    },
    button : {
        marginLeft: 5
    }
};

function generate(element) {
    return [0, 1, 2].map(value =>
      React.cloneElement(element, {
        key: value,
      }),
    );
  }

class Settings extends Component {

    constructor(props) {

        super(props);
        this.state = {
            protected: true,
            email: 'taicow@gmail.com',
            handle: 'taicow',
            password: '123',
        }
    }

    handleProtectedCheck = (e) => {

        this.setState({
            protected: false
        })
    }

    handleEmailChange = (e) => {
        
        console.log("something")
    }

    render () {
    return (
        <Fragment>
            <NavBar />
            <Grid style={styles.grid.container}  container spacing={24} >
                
                <Grid item xs={2}>
                </Grid>
                <Grid item xs={8}>


                            <Card>
                                <CardHeader
                                    title="Settings"
                                />
                                <hr />
                                <CardContent>
                                    <FormGroup row>
                                        <FormControlLabel
                                            control={
                                                <Checkbox
                                                checked={this.state.protected}
                                                onChange={this.handleProtectedCheck}
                                                name="protected"
                                                value="checked"
                                                />
                                            }
                                            label="Protected Account"
                                        />
                                    </FormGroup>

                                    <div>
                                        <TextField
                                            id="email"
                                            name="email"
                                            label="Email"
                                            margin="normal"
                                            value={this.state.email}
                                        />
                                        <Button type="submit" styles={styles.button}>
                                            Change Email
                                        </Button>
                                    </div>

                                    <div>
                                        <TextField
                                            id="handle"
                                            name="handle"
                                            label="Handle"
                                            margin="normal"
                                            value={this.state.handle}
                                        />
                                        <Button type="submit" styles={styles.button}>
                                            Change Handle
                                        </Button>
                                    </div>

                                    <div>
                                        <TextField
                                            id="password"
                                            name="password"
                                            label="Password"
                                            margin="normal"
                                            type="password"
                                            value={this.state.password}
                                        />
                                        <Button type="submit" styles={styles.button}>
                                            Change Password
                                        </Button>
                                    </div>
                                </CardContent>
                            </Card>


                </Grid>
                <Grid item xs={2}>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default withStyles(styles)(Settings);
