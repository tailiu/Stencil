import React, {Component, Fragment} from "react";

import TextField from 'material-ui/TextField';
import Grid from 'material-ui/Grid';
import Button from 'material-ui/Button';
import Card, {CardContent, CardHeader } from 'material-ui/Card';
import NavBar from './NavBar';
import Checkbox from 'material-ui/Checkbox';
import { FormGroup, FormControlLabel } from 'material-ui/Form';
import MessageBar from './MessageBar';
import axios from 'axios';
import { withCookies } from 'react-cookie';

const styles = {
    grid : {
        container : {
            marginTop: 80,
            height: 200
        }
    },
    button : {
        marginLeft: 5
    }
};

class Settings extends Component {

    constructor(props) {

        super(props);
        this.cookies = this.props.cookies;
        this.state = {
            user_id: this.cookies.get('user_id'),
            protected: false,
            email: '',
            handle: '',
            password: '',
        }
    }

    componentDidMount(){
        this.getUserInfo();
    }

    getUserInfo(){
        axios.get(
          'http://localhost:3000/users/getUserInfo',
          {
            params: {
              'user_id': this.state.user_id, 
              "req_token": this.cookies.get('req_token')
            }
          }
        ).then(response => {
          if(response.data.result.success){
            this.setState({
                protected: response.data.result.user.protected,
                email: response.data.result.email,
                handle: response.data.result.user.handle,
                // password: "123456"
            })
          }else{
            this.MessageBar.showSnackbar("User doesn't exist!");
          }
        })
    }

    handleProtectedCheck = (e) => {

        const isProtected = this.state.protected;

        this.setState({
            protected: !isProtected
        })
    }

    handleEmailChange = (email, e) => {
        
        axios.get(
            'http://localhost:3000/users/updateEmail',
            {
              params: {
                'user_id': this.state.user_id, 
                'email': this.state.email, 
              }
            }
          ).then(response => {
              console.log(response)
            if(response.data.result.success){
              this.setState({
                  email: response.data.result.user.email,
              })
              this.MessageBar.showSnackbar("Email changed!");
            }else{
              this.MessageBar.showSnackbar("Email can't be changed!");
            }
          })
    }

    handlePasswordChange = (e) => {
        
        axios.get(
            'http://localhost:3000/users/updatePassword',
            {
              params: {
                'user_id': this.state.user_id, 
                'password': this.state.password, 
                "req_token": this.cookies.get('req_token')
              }
            }
          ).then(response => {
              console.log(response)
            if(response.data.result.success){
              this.MessageBar.showSnackbar("Password changed!");
            }else{
              this.MessageBar.showSnackbar("Password can't be changed!");
            }
          })
    }

    handleHandleChange = (e) => {
        
        axios.get(
            'http://localhost:3000/users/updateHandle',
            {
              params: {
                'user_id': this.state.user_id, 
                'handle': this.state.handle, 
              }
            }
          ).then(response => {
              console.log(response)
            if(response.data.result.success){
              this.setState({
                  handle: response.data.result.user.handle,
              })
              this.MessageBar.showSnackbar("Handle changed!");
            }else{
              this.MessageBar.showSnackbar("Handle can't be changed!");
            }
          })
    }

    onChangeHandle =(e)=> {
        this.setState({
            "handle": e.target.value
        })
    }

    onChangeEmail =(e)=> {
        this.setState({
            "email": e.target.value
        })
    }

    onChangePassword =(e)=> {
        this.setState({
            "password": e.target.value
        })
    }

    onChangeProtected =(e)=> {

        axios.get(
            'http://localhost:3000/users/updateProtected',
            {
              params: {
                'user_id': this.state.user_id, 
                'protected': !this.state.protected, 
                "req_token": this.cookies.get('req_token')
              }
            }
          ).then(response => {
              console.log(response)
            if(response.data.result.success){
              this.setState({
                  protected: response.data.result.user.protected,
              })
              this.MessageBar.showSnackbar("Account protection changed!");
            }else{
              this.MessageBar.showSnackbar("Protected account can't be changed!");
            }
        })
    }

    render () {
    return (
        <Fragment>
            <NavBar />
            <Grid style={styles.grid.container}  container spacing={24} >
                
                <Grid item xs={2}>
                </Grid>
                <Grid item xs={8}>


                            <Card>
                                <CardHeader
                                    title="Settings"
                                />
                                <hr />
                                <CardContent>
                                    <FormGroup row>
                                        <FormControlLabel
                                            control={
                                                <Checkbox
                                                checked={this.state.protected}
                                                // onChange={this.handleProtectedCheck}
                                                onChange={this.onChangeProtected}
                                                name="protected"
                                                value="checked"
                                                />
                                            }
                                            label="Protected Account"
                                        />
                                    </FormGroup>

                                    <div>
                                        <TextField
                                            id="email"
                                            name="email"
                                            label="Email"
                                            margin="normal"
                                            value={this.state.email}
                                            onChange={this.onChangeEmail}
                                        />
                                        <Button type="submit" style={styles.button} onClick={this.handleEmailChange}>
                                            Change Email
                                        </Button>
                                    </div>

                                    <div>
                                        <TextField
                                            id="handle"
                                            name="handle"
                                            label="Handle"
                                            margin="normal"
                                            value={this.state.handle}
                                            onChange={this.onChangeHandle}
                                        />
                                        <Button type="submit" style={styles.button} onClick={this.handleHandleChange}>
                                            Change Handle
                                        </Button>
                                    </div>

                                    <div>
                                        <TextField
                                            id="password"
                                            name="password"
                                            label="Password"
                                            margin="normal"
                                            type="password"
                                            value={this.state.password}
                                            onChange={this.onChangePassword}
                                        />
                                        <Button type="submit" style={styles.button} onClick={this.handlePasswordChange}>
                                            Change Password
                                        </Button>
                                    </div>
                                </CardContent>
                            </Card>


                </Grid>
                <Grid item xs={2}>
                <MessageBar ref={instance => { this.MessageBar = instance; }}/>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default withCookies(Settings);
