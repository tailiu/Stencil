import React, {Component} from "react"
import PropTypes from "prop-types"

class Welcome extends Component {
	constructor(props) {
		super(props);

		this.handleSubmit = this.handleSubmit.bind(this);
	}
  
	handleSubmit(event) {
		event.preventDefault();
	}
  
	render() {		
		return (
			<div>
				<h1>See what’s happening in the world right now</h1>
				<h3>Join Twitter Today.</h3>
				<form action="http://localhost:3000/pages/index" onSubmit={this.handleSubmit}>
					<input type="submit" value="Sign Up" /><br/>
					<input type="submit" value="Login" />
				</form>
			</div>
	  	);
	}
  }

export default Welcome