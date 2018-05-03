import React, {Component, Fragment} from 'react';
import axios from 'axios';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
import Dialog, {
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
  } from 'material-ui/Dialog';
import Button from 'material-ui/Button';
import TextField from 'material-ui/TextField';
import MessageBar from './MessageBar';

class NewTweetDialog extends Component {

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
      };

    constructor(props) {
        
        super(props);
        
        const { cookies } = this.props;

        this.state = {
            user_id: cookies.get('user_id'),
            user_name: cookies.get('user_name'),
            user_handle: cookies.get('user_handle'),
            value : 0,
            anchorEl: null,
            tweet_box_open: this.props.open,
            tweet_content: "",
        }

    }

    componentWillReceiveProps(newProps) {
        this.setState({tweet_box_open: newProps.open});
    }

    handleTweetBoxClose = () => {
        this.setState({ tweet_box_open: false });
        this.props.onChange();
    };

    updateTweetContent = (e) => {
        this.setState({
            tweet_content: e.target.value
        })
    }

    validateForm = () => {
        if(this.state.tweet_content) 
            return true;
        else return false;
    }

    handleNewTweet = (e) => {
        
        if(!this.validateForm()){
          this.MessageBar.showSnackbar("Tweet box can't be empty!")
        }else{
            
          axios.get(
            'http://localhost:3000/tweets/new',
            {
              params: {
                'content':this.state.tweet_content, 
                'user_id': this.state.user_id,
                'reply_id': this.props.reply_id
              }
            }
          ).then(response => {
            console.log("axios:"+JSON.stringify(response))
            if(!response.data.result.success){
              this.MessageBar.showSnackbar(response.data.result.error.message)
            }else{
              this.MessageBar.showSnackbar("Tweet Posted!");
              this.handleTweetBoxClose();
            }
          })
        }
        e.preventDefault();
    }

    render(){
        return(
            <Fragment>
                <MessageBar ref={instance => { this.MessageBar = instance; }}/>
                <Dialog
                    open={this.state.tweet_box_open}
                    onClose={this.handleTweetBoxClose}
                    aria-labelledby="form-dialog-title"
                    >
                    <DialogTitle id="form-dialog-title">New Tweet</DialogTitle>
                    <DialogContent>
                        <DialogContentText>
                        {/* What's on your mind? */}
                        </DialogContentText>
                        <TextField
                        autoFocus
                        margin="dense"
                        id="tweet"
                        label="What's on your mind?"
                        type="email"
                        value={this.state.tweet_content}
                        onChange={this.updateTweetContent}
                        fullWidth
                        />
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={this.handleTweetBoxClose} color="primary">
                        Video/Photo
                        </Button>
                        <Button onClick={this.handleTweetBoxClose} color="primary">
                        Cancel
                        </Button>
                        <Button onClick={this.handleNewTweet} color="primary">
                        Tweet!
                        </Button>
                    </DialogActions>
                </Dialog>
            </Fragment>
        );
    }

}

export default withCookies(NewTweetDialog);