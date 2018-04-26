import React, {Component} from "react";
import axios from 'axios';

import axios from "axios"

import TextField from 'material-ui/TextField';
import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';

import Button from 'material-ui/Button';
import Card, { CardActions, CardContent } from 'material-ui/Card';
import Snackbar from 'material-ui/Snackbar';

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
      handle : '',
      name : ''
    }
  }

  validateForm = () => {

  }

  handleSubmit = e => {

    this.setState({
      name:e.target.name.value,
      email:e.target.email.value,
      password:e.target.password.value,
      handle:e.target.handle.value
    })

    if(!this.validateForm()){

    }else{
      axios.get(
        'http://localhost:3000/users/new',
        {
          params: {
            'name':this.state.name, 
            'email':this.state.email, 
            'handle':this.state.handle,
            'password': this.state.password
          }
        }
      ).then(response => {
        console.log(response)
        this.setState({username: response.data.name})
      })
    }

    e.preventDefault();
  }

  handleChange = e => {
    this.setState({ value: e.target.value });
  }

  goToLogin = e => {
		window.location = '/login';
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
                <form onSubmit={this.handleSubmit}>

                  <TextField
                    id="name"
                    label="Name"
                    margin="normal"
                    style={styles.card.input}
                    name="name"
                    value={this.state.name.value}
                    onChange={this.handleChange}
                  />
                  <br/>
                  <TextField
                    id="handle"
                    label="Handle"
                    margin="normal"
                    name="handle"
                    value={this.state.handle.value}
                    style={styles.card.input}
                    onChange={this.handleChange}
                  />
                  <br/>
                  <TextField
                    id="email"
                    label="Email"
                    margin="normal"
                    style={styles.card.input}
                    name="email"
                    value={this.state.email.value}
                    onChange={this.handleChange}
                  />
                  <br/>
                  <TextField
                    id="password"
                    label="Password"
                    margin="normal"
                    type="password"
                    style={styles.card.input}
                    name="password"
                    value={this.state.password.value}
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

export default SignUp;
