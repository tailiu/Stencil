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

const styles = {
    grid : {
        container : {
            marginTop: 80
        }
    },
    card: {
        card:{
            // minWidth: 400,
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
    tweet: {
        avatar: {

        },
        main_input: {

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

    like(e) {

    }

    retweet(e) {
        
    }

    reply(e) {
        
    }

  render () {
    return (
        <Fragment>
            <NavBar />
            <Grid style={styles.grid.container} container spacing={24} align="center">
                
                <Grid item xs={4}>
                    <Card style={styles.card.card} align="left">
                        <CardHeader
                            avatar={
                                <Avatar aria-label="Recipe" style={styles.tweet.avatar}>
                                TC
                                </Avatar>
                            }
                            title="Tai Cow"
                            subheader="Followers:49, Following:51, Tweets:90"
                        />
                    </Card>

                </Grid>

                <Grid item xs={8}>
                    <Grid container direction="column" align="left">
                        <Grid item>
                            <Card style={styles.card.card}>
                                <CardHeader
                                    avatar={
                                        <Avatar aria-label="Recipe" style={styles.tweet.avatar}>
                                        TC
                                        </Avatar>
                                    }
                                    title="Tai Cow"
                                    subheader="September 14, 2016"
                                />
                                <CardContent>
                                    <Typography component="p">
                                        This impressive paella is a perfect party dish and a fun meal to cook together with
                                    your guests. Add 1 cup of frozen peas along with the mussels, if you like.
                                    </Typography>
                                
                                </CardContent>
                                <CardActions>
                                    <Button size="small" onClick={this.like}>
                                        Like
                                    </Button>
                                    <Button size="small" onClick={this.retweet}>
                                        Retweet
                                    </Button>
                                    <Button size="small" onClick={this.reply}>
                                        Reply
                                    </Button>
                                </CardActions>
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
