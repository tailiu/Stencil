import React, {Component} from "react";
import axios from 'axios';

import TextField from 'material-ui/TextField';

import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';

import Button from 'material-ui/Button';
import Card, { CardActions, CardContent } from 'material-ui/Card';

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

class Login extends Component {

    constructor(props) {
        super(props);
        this.state = {
            email : '',
            password : '',
            name : '',
            value : '',
        }
    }

    getValidationState = () => {
    const length = this.state.value.length;
    if (length > 10) return 'success';
    else if (length > 5) return 'warning';
    else if (length > 0) return 'error';
    return null;
    }

    handleChange = (e) => {
    // this.setState({ value: e.target.value });
    }

    handleLogin = e => {
      console.log("Handle login!");

      axios.get(
        'http://localhost:3000/users/new',
        {
          params: {'name':'tai', 'email':'tai@cow.com', 'handle':'taicow'}
        }
      ).then(response => {
        console.log(response)
        this.setState({username: response.data.name})
      })

      e.preventDefault();
    }

    goToSignUp= (e) => {
        window.location = '/signup';
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
                <form onSubmit={this.handleLogin}>

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

export default Login;
