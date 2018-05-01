import React, {Component} from "react";
import axios from 'axios';
import { instanceOf } from 'prop-types';
import TextField from 'material-ui/TextField';

import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';

import Button from 'material-ui/Button';
import Card, { CardActions, CardContent } from 'material-ui/Card';

import TitleBar from './TitleBar';
import MessageBar from './MessageBar';

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

    goToHome = () => {
      this.props.history.push({pathname: '/home'});
    }
  
    handleLogin = () =>  {
      const { cookies } = this.props;
      if(!this.validateForm()){
        this.MessageBar.showSnackbar("Some fields are left empty!")
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
          // console.log(response)
          if(!response.data.result.success){
            this.MessageBar.showSnackbar(response.data.result.error.message)
          }else{
            this.MessageBar.showSnackbar("Login Successful!");
            cookies.set('user_id',  response.data.result.user.id);
            cookies.set('user_name', response.data.result.user.name);
            cookies.set('user_handle', response.data.result.user.handle);
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

    return (
      <Grid container style={styles.grid} spacing={24} align="center">

        <Grid item xs>
          <TitleBar />
					<MessageBar ref={instance => { this.MessageBar = instance; }}/>
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
