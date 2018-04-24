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

import SearchIcon from 'images/search_icon.svg';

import List, { ListItem, ListItemIcon, ListItemText } from 'material-ui/List';
import Divider from 'material-ui/Divider';

import TwitterLogo from 'images/Twitter_Logo_Blue.png';

const styles = {
    grid : {
        container : {
            marginTop: 80
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

class Home extends Component {

    constructor(props) {

        super(props);
        this.state = {
        }
    }

  render () {
    return (
        <Fragment>
            <NavBar />
            <Grid style={styles.grid.container} direction="column" container spacing={24} align="center">

                <Grid item xs={8}>
                    <Grid container direction="column" align="left">
                        <Grid item>
                            <Card>
                                <CardHeader
                                    title="Search"
                                />
                                <hr />
                                <CardContent>

                                    <List dense={true}>
                                        {generate(
                                        <ListItem>
                                            <ListItemIcon>
                                            <img style={styles.logo} src={TwitterLogo} /> 
                                            {/* <FolderIcon /> */}
                                            </ListItemIcon>
                                            <ListItemText
                                            primary="Single-line item"
                                            //   secondary={secondary ? 'Secondary text' : null}
                                            />
                                        </ListItem>,
                                        )}
                                    </List>
                                </CardContent>
                            </Card>
                        </Grid>
                    </Grid>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default withStyles(styles)(Home);
