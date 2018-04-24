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

import TitleBar from './TitleBar';

const styles = {
  grid : {
    // background: "#c0deed",
    height: "100%"
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
      name : ''
    }

    this.handleSubmit = this.handleSubmit.bind(this);
    this.getValidationState = this.getValidationState.bind(this);
    this.handleChange = this.handleChange.bind(this);
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

  goToLogin(e) {
		window.location = 'http://localhost:3000/pages/login';
  }

  render () {
    return (
      <Grid container style={styles.grid} spacing={24} >

        <Grid item xs>
          <TitleBar />
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
                <form action="/users/new" onSubmit={this.handleSubmit}>

                  <TextField
                    id="name"
                    label="Name"
                    margin="normal"
                    style={styles.card.input}
                    onChange={this.handleChange}
                  />
                  <br/>
                  <TextField
                    id="handle"
                    label="Handle"
                    margin="normal"
                    style={styles.card.input}
                    onChange={this.handleChange}
                  />
                  <br/>
                  <TextField
                    id="email"
                    label="Email"
                    margin="normal"
                    style={styles.card.input}
                    onChange={this.handleChange}
                  />
                  <br/>
                  <TextField
                    id="password"
                    label="Password"
                    margin="normal"
                    type="password"
                    style={styles.card.input}
                    onChange={this.handleChange}
                  />
                  <br/>
                  <br/>
                  <Button type="submit" variant="raised" style={styles.card.button}>
                    Sign Up
                  </Button>
                </form>
              </CardContent>
              <CardActions>
                <Button size="small" onClick={this.goToLogin}>
                Already signed up? Login!
                </Button>
              </CardActions>
            </Card>
          </Grid>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles()(SignUp);
