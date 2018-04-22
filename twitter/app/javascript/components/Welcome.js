import React, {Component} from "react"
import PropTypes from "prop-types"

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
			<div>
				<h1>See what’s happening in the world right now</h1>
				<h3>Join Twitter Today.</h3>
				<input type="button" value="Login" onClick={this.handleClick}  />
				<input type="button" value="Sign Up" onClick={this.handleClick} />
			</div>
		  );
	}
  }

export default Welcome