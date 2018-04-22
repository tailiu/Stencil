import React, {Component, Fragment} from "react";

import { Link } from 'react-router';

import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import MenuItem from 'material-ui/Menu/MenuItem';
import TextField from 'material-ui/TextField';

import Paper from 'material-ui/Paper';
import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';

import AppBar from 'material-ui/AppBar';
import Toolbar from 'material-ui/Toolbar';

import Button from 'material-ui/Button';
import Card, { CardActions, CardContent, CardHeader } from 'material-ui/Card';

const styles = {
  grid : {
    // background: "#c0deed",
    height: "100%"
  },
  navbar : {
    navbar: {
      backgroundColor: "#00aced",
    },
    title: {
      color: "#fff",
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
  }
};

class SignUp extends Component {

  constructor(props) {
    super(props);
    this.state = {
      email : '',
      password : '',
      name : '',
      value : '',
    }
  }

  handleSignUp(event) {
    console.log("sign me up, bozo!");
    event.preventDefault();
  }

  getValidationState() {
    const length = this.state.value.length;
    if (length > 10) return 'success';
    else if (length > 5) return 'warning';
    else if (length > 0) return 'error';
    return null;
  }

  handleChange(e) {
    // this.setState({ value: e.target.value });
  }

  render () {
    return (
      <Grid container style={styles.grid} spacing={24} >

        <Grid item xs>
          <AppBar style={styles.navbar.navbar} position="static" color="default">
            <Toolbar>
              <Typography variant="title" style={styles.navbar.title}>
                Twitter
              </Typography>
            </Toolbar>
          </AppBar>
        </Grid>

        <Grid item xs={12}>
        </Grid>
        
        <Grid container spacing={24} direction="column" align="center">
          <Grid item xs>
          <Typography variant="headline" gutterBottom>
            <strong>Join Twitter Today!</strong>
          </Typography>
          </Grid>

          <Grid item xs>
            <Card style={styles.card.card}>
              {/* <CardHeader
                title="Join Twitter Today!"
              /> 
              <hr/> */}
              <CardContent>
                <form onSubmit={this.handleSignUp}>

                  <TextField
                    id="name"
                    label="Name"
                    margin="normal"
                    style={styles.card.input}
                  />
                  <br/>
                  <TextField
                    id="email"
                    label="Email"
                    margin="normal"
                    style={styles.card.input}
                  />
                  <br/>
                  <TextField
                    id="password"
                    label="Password"
                    margin="normal"
                    type="password"
                    style={styles.card.input}
                  />
                  <br/>
                  <br/>
                  <Button type="submit" variant="raised" style={styles.card.button}>
                    Sign Up
                  </Button>
                </form>
              </CardContent>
              <CardActions>
                <Button size="small">Already signed up? Login!</Button>
              </CardActions>
            </Card>
          </Grid>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles()(SignUp);
