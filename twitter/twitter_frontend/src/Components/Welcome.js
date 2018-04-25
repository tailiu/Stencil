import React, {Component} from "react";
import PropTypes from "prop-types";
import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';
import Button from 'material-ui/Button';
import Card, { CardActions, CardContent, CardHeader } from 'material-ui/Card';


const styles = {
	
	logo: {
		height: 150,
	},
	button: {
		// width: "100%",
		backgroundColor: "#00aced",
		color: "#fff",
		margin: 5
	  }
  };

class Welcome extends Component {
	constructor(props) {
		super(props);

		this.handleClick = this.handleClick.bind(this);
	}

	handleClick(val, e) {

		if (val.toLowerCase() == "login") {
			window.location = '/login';
		} else {
			window.location = '/signUp';
		}
	}

	render() {		
		return (
			<Grid container spacing={24} direction="column" align="center">
				
				<Grid item xs>
					<img style={styles.logo} src={require('../Assets/Images/Twitter_Logo_Blue.png')} /> 
				</Grid>
				
				<Grid item xs>
					<Typography variant="display1" gutterBottom>
							<strong>see what’s happening in the world right now</strong>
					</Typography>
				</Grid>

				<Grid item xs>
				</Grid>

				<Grid item xs>
					<Typography variant="headline" gutterBottom>
						<strong>Join Twitter Today!</strong>
					</Typography>
				</Grid>

				<Grid item xs>

					<Button 
						style={styles.button}
						type="submit" 
						variant="raised" 
						value="Login" 
						onClick={this.handleClick.bind(this, "login")}>
						Login
					</Button>

					<Button 
						style={styles.button}
						type="submit" 
						variant="raised" 
						value="Sign Up" 
						onClick={this.handleClick.bind(this, "signup")}>
						Sign Up
					</Button>
				</Grid>

			</Grid>
	  	);
	}
  }

export default Welcome