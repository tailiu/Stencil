import React, {Component} from 'react';
import Snackbar from 'material-ui/Snackbar';

class MessageBar extends Component {

    constructor(props){
        super(props);
        this.state = {
            show: false,
            message: "Some Message!"
        }
    }

    showSnackbar = message => {
        console.log("HERE SNACKS!");
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
                anchorOrigin={{
                vertical: 'top',
                horizontal: 'center',
                }}
                open={this.state.show}
                autoHideDuration={6000}
                // onClose={this.handleClose}
                SnackbarContentProps={{
                'aria-describedby': 'message-id',
                }}
                message={<span id="message-id">{this.state.message}</span>}
                action={[
                ]}
            />
        )
    }
}

export default MessageBar;