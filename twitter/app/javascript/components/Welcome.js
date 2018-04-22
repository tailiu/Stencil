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

		this.handleClick = this.handleClick.bind(this);
	}

	handleClick(event) {
		if (event.target.value == "Login") {
			window.location = 'http://localhost:3000/pages/login';
		} else {
			window.location = 'http://localhost:3000/pages/signUp';
		}
	}

	render() {		
		return (
			<Grid container style={styles.grid} spacing={24} >
			<div>
				<h1>See what’s happening in the world right now</h1>
				<h3>Join Twitter Today.</h3>
				<input type="button" value="Login" onClick={this.handleClick}  />
				<input type="button" value="Sign Up" onClick={this.handleClick} />
			</div>
			</Grid>
	  	);
	}
  }

export default Welcome