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

import FavoriteIcon from 'images/Twitter_Logo_Blue.png';
import ShareIcon from 'images/Twitter_Logo_Blue.png';
import ExpandMoreIcon from 'images/Twitter_Logo_Blue.png';
import MoreVertIcon from 'images/Twitter_Logo_Blue.png';

const styles = {
    grid : {
        container : {
            marginTop: 50
        }
    },
    card: {
        card:{
            minWidth: 375,
        },
        input:{
            width: "95%",
        },
        button: {
            width: "100%",
            backgroundColor: "#00aced",
            color: "#fff",
        }
    },
    paper: {
        height: "100%",
        width: "100%",
        // margin: 20,
        textAlign: 'center',
        display: 'inline-block',
    },
    tweet: {
        main_input: {
            minWidth: 400       
        },
        container: {

        },
        card: {

        }
    }
};

class Home extends Component {

  constructor(props) {
    
    super(props);

    this.state = {
      email : '',
      password : '',
      name : '',
      tweet_value : '',
      value : '',
    }

    this.handleSubmit = this.handleSubmit.bind(this);
    this.getValidationState = this.getValidationState.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleChangeToTweetField = this.handleChangeToTweetField.bind(this);
    this.goToLogin = this.goToLogin.bind(this);
  }

  handleSubmit(e) {
    console.log("Called: 'handleSignUp'");

    // this.state.name = e.target.name.value;
    // this.state.email = e.target.email.value;
    // this.state.password = e.target.password.value;
    
    alert(this.state.name);
    e.preventDefault();
  }

  getValidationState() {
    const length = this.state.value.length;
    if (length > 10) return 'success';
    else if (length > 5) return 'warning';
    else if (length > 0) return 'error';
    return null;
  }

  handleChange(e) {
    this.setState({ value: e.target.value });
  }

  handleChangeToTweetField(e) {
    this.setState({ tweet_value: e.target.value });
  }

  goToLogin(e) {
		window.location = 'http://localhost:3000/pages/login';
  }

  render () {
    return (
        <Fragment>
            <NavBar />
            <Grid style={styles.grid.container} container spacing={24} direction="column" align="center">

                <Grid item xs>
                    
                </Grid>
                
                <Grid item xs>

                    <TextField
                    id="tweet"
                    label="  What's on your mind?"
                    margin="normal"
                    style={styles.tweet.main_input}
                    onChange={this.handleChangeToTweetField}
                  />

                </Grid>

                <Grid item xs>
                    <Card style={styles.card.card}>
                        {/* <CardHeader
                        title="Join Twitter Today!"
                        /> 
                        <hr/> */}
                        <CardContent>
                        
                        </CardContent>
                        <CardActions>
                            <Button size="small" onClick={this.goToLogin}>
                                something
                            </Button>
                        </CardActions>
                    </Card>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default withStyles(styles)(Home);
