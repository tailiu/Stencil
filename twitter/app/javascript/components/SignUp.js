import React, {Component, Fragment} from "react";

import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import MenuItem from 'material-ui/Menu/MenuItem';
import TextField from 'material-ui/TextField';

import Paper from 'material-ui/Paper';
import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';

import AppBar from 'material-ui/AppBar';
import Toolbar from 'material-ui/Toolbar';

const styles = {
  "background" : {
    backgroundColor: "#c0deed"
  },
  "navbar" : {
    backgroundColor: "#00aced",
    color: "#fff"
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

      
      <Grid styles={styles.background}>

        <Grid item xs>
          <AppBar style={styles.navbar} position="static" color="default">
            <Toolbar>
              <Typography variant="title" >
                Twitter
              </Typography>
            </Toolbar>
          </AppBar>
        </Grid>
        
        <Grid container spacing={24} direction="column" align="center">
          <Grid item xs>
          <Typography variant="headline" gutterBottom>
            Join Twitter Today!
          </Typography>
          </Grid>

          <Grid item xs>
          <form onSubmit={this.handleSignUp}>

            <TextField
              id="name"
              label="Name"
              margin="normal"
            />
            <br/>
            <TextField
              id="email"
              label="Email"
              margin="normal"
            />
            <br/>
            <TextField
              id="password"
              label="Password"
              margin="normal"
            />

          </form>
          </Grid>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles()(SignUp);
