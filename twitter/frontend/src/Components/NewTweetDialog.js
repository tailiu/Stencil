import React, {Component, Fragment} from 'react';
import axios from 'axios';
import { withCookies } from 'react-cookie';
import Dialog, {
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
  } from 'material-ui/Dialog';
import Button from 'material-ui/Button';
import TextField from 'material-ui/TextField';
import MessageBar from './MessageBar';
import Card from 'material-ui/Card';
import { CardContent } from 'material-ui';

axios.defaults.withCredentials = true

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

    constructor(props) {
        
        super(props);
        
        this.cookies = this.props.cookies;

        this.state = {
            user_id: this.cookies.get('user_id'),
            user_name: this.cookies.get('user_name'),
            user_handle: this.cookies.get('user_handle'),
            value : 0,
            anchorEl: null,
            tweet_box_open: this.props.open,
            tweet_content: "",
            hasMedia: false,
            mediaUrl: '../Assets/Images/liked-icon.png',
            imagePreviewUrl: '',
            file: '',
            media_type: 'text',
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

        let file = false
        if (this.state.hasMedia) file = this.state.file
        
        let formData = new FormData();

        formData.set('content', this.state.tweet_content);
        formData.set('user_id',this.state.user_id);
        formData.set('type', this.state.media_type);
        formData.set('file',file);
        formData.set('reply_id', this.props.reply_id);
        formData.set('req_token', this.cookies.get('req_token'));

        axios('http://localhost:8000/tweets/newf', {
            method: 'post',
            data: formData,
            withCredentials: true,
            config: { 
                headers: {
                    'Content-Type': 'multipart/form-data', 
                }
            }
        })
        .then((response) => {
            console.log(response.data.result);
            if(!response.data.result.success){
                this.MessageBar.showSnackbar(response.data.result.error.message)
            }else{
                this.MessageBar.showSnackbar("Tweet Posted!");
                this.handleTweetBoxClose();
            }
        })
        // .catch(function (response) {
        //     console.log("something wrong ")
        //     console.log(response);
        // });
    }

    handleImageChange =(e)=> {
        e.preventDefault();

        if (e.target.files.length <= 0) {
            console.log("No files selected.")
            return
        }

        if(e.target.files[0].type.indexOf("image") >= 0){
            this.setState({
                "media_type":"photo"
            })
        } else if (e.target.files[0].type.indexOf("video") >= 0) {
            this.setState({
                "media_type":"video"
            })
        } else {
            this.MessageBar.showSnackbar("This type of file is not allowed!");
            return false;
        }
    
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
                        {this.state.hasMedia &&
                            <Card style={styles.upload.area}>
                                {this.state.media_type === "photo" &&
                                    <CardContent>
                                        <img alt="preview" style={styles.upload.image} src={this.state.imagePreviewUrl} />
                                    </CardContent>
                                }
                                {this.state.media_type === "video" &&
                                    <CardContent>
                                        <video width="320" height="240" controls>
                                            <source src={this.state.imagePreviewUrl} type="video/mp4"/>
                                            Your browser does not support the video tag.
                                        </video>
                                    </CardContent>
                                }
                            </Card>
                        }
                    </DialogContent>
                    <DialogActions>

                        <Button type="file " color="primary" label='Video/Photo'>
                            <input type="file" style={styles.upload.input} onChange={this.handleImageChange} />
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