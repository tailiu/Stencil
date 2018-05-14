import React, {Component, Fragment} from 'react';
import axios, {post} from 'axios';
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
            imagePreviewUrl: '',
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

    fileUpload(){
        const url = 'http://localhost:3000/tweets/newf';
        const formData = new FormData();
        let file = false
        if (this.state.hasMedia){
            file = this.state.file
        }
        formData.append('content', this.state.tweet_content);
        formData.append('user_id',this.state.user_id);
        formData.append('type', this.state.media_type);
        formData.append('file',file);
        formData.append('reply_id', this.props.reply_id);
        const config = {
            headers: {
                'content-type': 'multipart/form-data'
            }
        }
        return  post(url, formData, config)
    }    

    handleNewTweet = (e) => {
        
        if(!this.validateForm()){
          this.MessageBar.showSnackbar("Tweet box can't be empty!")
        }else{
            this.fileUpload(e).then((response)=>{
                if(!response.data.result.success){
                    this.MessageBar.showSnackbar(response.data.result.error.message)
                }else{
                    this.MessageBar.showSnackbar("Tweet Posted!");
                    this.handleTweetBoxClose();
                }
            })
            
        //   axios.get(
        //     'http://localhost:3000/tweets/new',
        //     {
        //       params: {
        //         'content':this.state.tweet_content, 
        //         'user_id': this.state.user_id,
        //         'reply_id': this.props.reply_id,
        //         'file': file,
        //       }
        //     }
        //   ).then(response => {
        //     console.log("axios:"+JSON.stringify(response))
        //     if(!response.data.result.success){
        //       this.MessageBar.showSnackbar(response.data.result.error.message)
        //     }else{
        //       this.MessageBar.showSnackbar("Tweet Posted!");
        //       this.handleTweetBoxClose();
        //     }
        //   })
        }
        e.preventDefault();
    }
    
    _handleImageChange =(e)=> {
        e.preventDefault();

        if(e.target.files[0].type.indexOf("image") >= 0){
            this.state.media_type = "photo"
        } else if (e.target.files[0].type.indexOf("video") >= 0) {
            this.state.media_type = "video"
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
                                        <img style={styles.upload.image} src={this.state.imagePreviewUrl} />
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