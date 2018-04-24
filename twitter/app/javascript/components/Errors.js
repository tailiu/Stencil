import React from "react"
import PropTypes from "prop-types"
class Errors extends React.Component {
  render () {
    return (
      <React.Fragment>
        <h1>Errors: {this.props.errors}</h1>
      </React.Fragment>
    );
  }
}

Errors.propTypes = {
  errors: PropTypes.string
};
export default Errors
