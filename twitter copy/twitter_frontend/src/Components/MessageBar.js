import React, {Component} from 'react';
import Snackbar from 'material-ui/Snackbar';

const styles = {
    message: {
        textAlign: "center"
    },
    MessageBar: {
        zIndex: 1500
    }
}

class MessageBar extends Component {

    constructor(props){
        super(props);
        this.state = {
            show: false,
            message: "Some Message!"
        }
    }

    showSnackbar = message => {
        this.setState({
            message: message,
            show: true
        })
        setTimeout(function() { 
          this.setState({
              message: "",
              show: false
          }); 
        }.bind(this), 5000);
      }

    render(){
        return(
            <Snackbar
                style={styles.MessageBar}
                anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'center',
                }}
                open={this.state.show}
                autoHideDuration={2000}
                // onClose={this.handleClose}
                SnackbarContentProps={{
                    'aria-describedby': 'message-id',
                }}
                message={<span id="message-id" style={styles.message}>{this.state.message}</span>}
                action={[
                ]}
            />
        )
    }
}

export default MessageBar;