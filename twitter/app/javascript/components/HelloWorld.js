import React from "react"
import PropTypes from "prop-types"
import TaiCow from './TaiCow.js'
import SignUp from './SignUp.js'

const styles = {
  "background" : {
    backgroundColor: "#c0deed"
  }
}

class HelloWorld extends React.Component {

  handleClick() {
    console.log("something here");

  }

  render () {
    return (
      <div style={styles.background} >
        <SignUp />

      </div>
    );
  }
}

HelloWorld.propTypes = {
  greeting: PropTypes.string
};
export default HelloWorld
