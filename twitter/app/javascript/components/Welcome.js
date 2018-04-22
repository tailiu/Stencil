import React, {Component} from "react"
import PropTypes from "prop-types"

const styles = {
	grid : {
	  // background: "#c0deed",
	  height: "100%"
	},
	navbar : {
	  navbar: {
		backgroundColor: "#00aced",
	  },
	  title: {
		color: "#fff",
	  }
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

class Welcome extends Component {
	constructor(props) {
		super(props);

		this.state = {value: ''};

		this.handleClick = this.handleClick.bind(this);
	}

	handleClick(event) {
		this.setState({value: event.target.value});
	}

	render() {		
		return (
			<Grid container style={styles.grid} spacing={24} >
			<div>
				<h1>See what’s happening in the world right now</h1>
				<h3>Join Twitter Today.</h3>
				<form action="http://localhost:3000/pages/loginOrSignUp">
					<input type="hidden" name="userAction" value={this.state.value} /><br/>
					<input type="submit" value="Sign Up" onClick={this.handleClick} /><br/>
					<input type="submit" value="Login" onClick={this.handleClick} />
				</form>
			</div>
			</Grid>
	  	);
	}
  }

export default Welcome