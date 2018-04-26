import React, {Component} from "react";

import axios from 'axios';

import TextField from 'material-ui/TextField';
import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';

import Button from 'material-ui/Button';
import Card, { CardActions, CardContent } from 'material-ui/Card';
import Snackbar from 'material-ui/Snackbar';

const styles = {
  logo: {
		height: 150,
	},
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
      name : '',
      snackbar: {
        message: 'Some Error!',
        show: false
      }
    }
  }

  validateForm = () => {
    if(this.state.name && 
      this.state.email &&
      this.state.password &&
      this.state.handle
    ) return true;
    else return false;
  }

  showSnackbar = message => {
    this.setState({
      snackbar: {
        message: message,
        show: true
      }
    })
    setTimeout(function() { 
      this.setState({
        snackbar: {
          message: "",
          show: false
        }
      }); 
    }.bind(this), 5000);
  }

  handleSignup = () =>  {
    if(!this.validateForm()){
      this.showSnackbar("Some fields are left empty!")
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
        if(!response.data.result.success){
          this.showSnackbar(response.data.result.error.message)
        }else{
          this.setState({username: response.data.result.user.name})
          this.showSnackbar("Signup Successful!");
        }
      })
    }
  }

  handleSubmit = e => {

    this.setState({
      name:e.target.name.value,
      email:e.target.email.value,
      password:e.target.password.value,
      handle:e.target.handle.value
    }, () => this.handleSignup())

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
      <Grid container style={styles.grid} spacing={24} align="center">

        <Grid item xs>
					<img style={styles.logo} alt="Logo" src={require('../Assets/Images/Twitter_Logo_Blue.png')} /> 
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
                <Snackbar
                  anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'center',
                  }}
                  open={this.state.snackbar.show}
                  autoHideDuration={6000}
                  // onClose={this.handleClose}
                  SnackbarContentProps={{
                    'aria-describedby': 'message-id',
                  }}
                  message={<span id="message-id">{this.state.snackbar.message}</span>}
                  action={[
                  ]}
                />
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
