import React, {Component, Fragment} from "react";

import TextField from 'material-ui/TextField';
import Grid from 'material-ui/Grid';

import Button from 'material-ui/Button';
import Card, {CardContent, CardHeader } from 'material-ui/Card';

import NavBar from './NavBar';

import Checkbox from 'material-ui/Checkbox';

import { FormGroup, FormControlLabel } from 'material-ui/Form';

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
        this.state = {
            protected: true,
            email: 'taicow@gmail.com',
            handle: 'taicow',
            password: '123',
        }
    }

    handleProtectedCheck = (e) => {

        const isProtected = this.state.protected;

        this.setState({
            protected: !isProtected
        })
    }

    handleEmailChange = (e) => {
        
        console.log("Change Email")
    }

    handlePasswordChange = (e) => {
        
        console.log("Change Password")
    }

    handleHandleChange = (e) => {
        
        console.log("Change Handle")
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
                                                onChange={this.handleProtectedCheck}
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
                                        />
                                        <Button type="submit" style={styles.button} onClick={this.handlePasswordChange}>
                                            Change Password
                                        </Button>
                                    </div>
                                </CardContent>
                            </Card>


                </Grid>
                <Grid item xs={2}>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default Settings;
