import React, {Component} from "react";
import axios from 'axios';
import { instanceOf } from 'prop-types';
import TextField from 'material-ui/TextField';

import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';

import Button from 'material-ui/Button';
import Card, { CardActions, CardContent } from 'material-ui/Card';

import Snackbar from 'material-ui/Snackbar';

import { withCookies, Cookies } from 'react-cookie';

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

// const cookies = new Cookies(req.headers.cookie);

class Login extends Component {

  static propTypes = {
    cookies: instanceOf(Cookies).isRequired
  };

    constructor(props) {
        super(props);
        this.state = {
            email : '',
            password : '',
            snackbar: {
              message: 'Some Error!',
              show: false
            }
        }
    }

    validateForm = () => {
      if(this.state.email &&
        this.state.password
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

    goToHome = (e) => {
      window.location = '/home';
    }
  
    handleLogin = () =>  {
      const { cookies } = this.props;
      if(!this.validateForm()){
        this.showSnackbar("Some fields are left empty!")
      }else{
        axios.get(
          'http://localhost:3000/users/verify',
          {
            params: {
              'email':this.state.email, 
              'password': this.state.password
            }
          }
        ).then(response => {
          console.log(response)
          if(!response.data.result.success){
            this.showSnackbar(response.data.result.error.message)
          }else{
            this.showSnackbar("Login Successful!");
            cookies.set('session_id', response.data.result.session_id);
            setTimeout(function() { 
              this.goToHome();
            }.bind(this), 1000);
          }
        })
      }
    }
  
    handleSubmit = e => {
  
      this.setState({
        email:e.target.email.value,
        password:e.target.password.value,
      }, () => this.handleLogin())
  
      e.preventDefault();
    }

    goToSignUp = (e) => {
        window.location = '/signup';
    }

  render () {
    const { cookies } = this.props;
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
          <strong>Log In To Twitter!</strong>
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
                    Log In
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
                <Button size="small" onClick={this.goToSignUp} >New to Twitter? Sign Up!</Button>
              </CardActions>
            </Card>
          </Grid>
        </Grid>
      </Grid>
    );
  }
}

export default withCookies(Login);
