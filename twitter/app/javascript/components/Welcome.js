import React, {Component} from "react"
import PropTypes from "prop-types"

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
			<div>
				<h1>See what’s happening in the world right now</h1>
				<h3>Join Twitter Today.</h3>
				<form action="/pages/loginOrSignUp">
					<input type="hidden" name="userAction" value={this.state.value} /><br/>
					<input type="submit" value="Sign Up" onClick={this.handleClick} /><br/>
					<input type="submit" value="Login" onClick={this.handleClick} />
				</form>
			</div>
	  	);
	}
  }

export default Welcome