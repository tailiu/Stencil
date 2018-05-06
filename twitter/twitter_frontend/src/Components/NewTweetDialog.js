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
import Upload from 'material-ui-upload/Upload';
import Card, { CardMedia } from 'material-ui/Card';
import { CardContent } from 'material-ui';

const styles = {
    upload: {
        button: {
            border: "1px solid #ccc",
            display: "inline-block",
            padding: "6px 12px",
            cursor: "pointer"

        },
        input: {
            display:"hidden"
        },
        area: {
            alignItems: "center",
            textAlign: "center",
            marginTop: 10,
        },
        image: {
            height: 250
        }
    }
}

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
            file: '',
            imagePreviewUrl: '',
            hasMedia: false,
            mediaUrl: '../Assets/Images/liked-icon.png',
            imagePreviewUrl: '',
            file: '',
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

    onFileLoad = (e, file) => {
        console.log(e.target.result, file.name);
    }

    _handleSubmit =(e)=> {
        e.preventDefault();
        // TODO: do something with -> this.state.file
        console.log('handle uploading-', this.state.file);
      }
    
    _handleImageChange =(e)=> {
        e.preventDefault();
    
        let reader = new FileReader();
        let file = e.target.files[0];
    
        reader.onloadend = () => {
          this.setState({
            file: file,
            imagePreviewUrl: reader.result,
            hasMedia: true
          });
        }
    
        reader.readAsDataURL(file)
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
                        {this.state.hasMedia?
                            <Card style={styles.upload.area}>
                                <CardContent>
                                    <img style={styles.upload.image} src={this.state.imagePreviewUrl} />
                                </CardContent>
                            </Card>
                        :
                        <div></div>
                        }
                    </DialogContent>
                    <DialogActions>

                        <Button type="file " color="primary" label='Video/Photo'>
                            <input type="file" style={styles.upload.input} onChange={this._handleImageChange} />
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