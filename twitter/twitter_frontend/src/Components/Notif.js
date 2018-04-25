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

import List, { ListItem, ListItemIcon, ListItemText } from 'material-ui/List';
import Divider from 'material-ui/Divider';

const styles = {
    grid : {
        container : {
            marginTop: 80,
            height: 200
        }
    },
};

function generate(element) {
    return [0, 1, 2].map(value =>
      React.cloneElement(element, {
        key: value,
      }),
    );
  }

class Notif extends Component {

    constructor(props) {

        super(props);
        this.state = {
        }
    }

  render () {
    return (
        <Fragment>
            <NavBar />
            <Grid style={styles.grid.container} direction="column" container spacing={24} >
                
                <Grid item xs={2}>
                </Grid>
                <Grid item xs={8}>


                            <Card>
                                <CardHeader
                                    title="Notifications"
                                />
                                <hr />
                                <CardContent>

                                    <List dense={true}>
                                        <ListItem>
                                            <ListItemIcon>
                                            <img style={styles.logo} src={RetweetIcon} /> 
                                            {/* <FolderIcon /> */}
                                            </ListItemIcon>
                                            <ListItemText
                                            primary="Tai retweeted your tweet"
                                            //   secondary={secondary ? 'Secondary text' : null}
                                            />
                                        </ListItem>
                                        <ListItem>
                                            <ListItemIcon>
                                            <img style={styles.logo} src={FavIcon} /> 
                                            {/* <FolderIcon /> */}
                                            </ListItemIcon>
                                            <ListItemText
                                            primary="Miro favorited your tweet"
                                            //   secondary={secondary ? 'Secondary text' : null}
                                            />
                                        </ListItem>
                                        <ListItem>
                                            <ListItemIcon>
                                            <img style={styles.logo} src={FollowIcon} /> 
                                            {/* <FolderIcon /> */}
                                            </ListItemIcon>
                                            <ListItemText
                                            primary="Major Tom followed you"
                                            //   secondary={secondary ? 'Secondary text' : null}
                                            />
                                        </ListItem>
                                    </List>
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

export default withStyles(styles)(Notif);
