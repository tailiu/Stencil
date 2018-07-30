import React, {Component} from "react";
import axios, {post} from 'axios';
import TextField from 'material-ui/TextField';
import Button from 'material-ui/Button';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
import Grid from 'material-ui/Grid';
import Card, { CardActions, CardContent } from 'material-ui/Card';
import FileUpload from '@material-ui/icons/FileUpload';

var styles = {
    inputContainer: {
        height: '100%'
    },
    messagesInput: {
        width: "75%",
        marginLeft: 5,
        marginRight: 10
    },
    upload: {
        width: "100%",
    },
    preview: {
        textAlign: "center",
        marginTop: 5
    },
    media_preview: {
        height: 200,
        width: "auto"
    }
}

class MessageInputAllow extends Component {
    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    }

    constructor(props) {
        super(props);

        const { cookies } = this.props;

        this.state = {
            value: '',
            user_id: cookies.get('user_id'),
            media_type: null,
            media_preview: '',
            file: '',
            has_media: false
        }

        this.setFileName = React.createRef();
    }

    handleChange = (e) => {
       this.setState({
           value: e.target.value
       })
    }

    handleNewMessage = () => {
        var formData = new FormData();
        formData.append('content', this.state.value);
        formData.append('user_id',this.state.user_id);
        formData.append("conversation_id", this.props.current_conversation_id);
        formData.append('media_type', this.state.media_type);
        formData.append('media', this.state.file);
        const config = {
            headers: {
                'content-type': 'multipart/form-data'
            }
        };
        const url = 'http://localhost:3000/messages/new';

        post(url, formData, config).then(response => {
            if(!response.data.result.success) {
                this.props.messageBar.showSnackbar(response.data.result.error)
            }else{
                this.setState({
                    value: '',
                    media_type: null,
                    media_preview: '',
                    has_media: false,
                    file: ''
                })
                this.setFileName.current.value = ''
                this.props.setHasMediaState(false)
                this.props.onNewMessage(response.data.result.newMessage)
            }
        })
    }

    catchReturn = (e) => {
        if (e.key === 'Enter' && (this.state.value != '' || this.state.has_media) && this.props.current_conversation_id != '') {
            this.handleNewMessage()
        }
    }
    
    setMediaType = (media_type) => {
        this.setState({
            "media_type": media_type
        })
    }

    handleUploadChange = (e) => {
        const file = e.target.files[0]

        if (file == undefined) {
            return 
        }

        if (file.type.indexOf("image") < 0 && file.type.indexOf("video") < 0) {
            this.props.messageBar.showSnackbar("File type is not supported")
            return
        }

        const reader = new FileReader();
        reader.readAsDataURL(file)
        reader.onloadend = () => {
            if (file.type.indexOf("image") >= 0) {
                this.setMediaType("photo")
            } else if (file.type.indexOf("video") >= 0) {
                this.setMediaType("video")
            }
            this.setState({
                file: file,
                media_preview: reader.result,
                has_media: true
            })
        }

        this.props.setHasMediaState(true)
    }
    
    renderMedia = () => {
        if (this.state.has_media) {
            if (this.state.media_type == "photo") {
                return (
                    <CardContent style={styles.preview}>
                        <img style={styles.media_preview} src={this.state.media_preview} />
                    </CardContent>
                ) 
            } else if (this.state.media_type == "video") {
                return (
                    <CardContent style={styles.preview}>
                        <video style={styles.media_preview} controls>
                            <source src={this.state.media_preview} type="video/mp4"/>
                        </video>
                    </CardContent>
                )
            }
        }
    }

    render() {
        var disabled = true

        if ((this.state.value != '' || this.state.has_media) && this.props.current_conversation_id != '') {
            disabled = false
        }
        
        styles.sendMessageButton = {
            backgroundColor: disabled ? '#BBDEFB' : "#00aced",
            color: "#fff",
            variant: "raised"
        }

        return (
            <Card style={styles.inputContainer}>
                {this.renderMedia()}
                <CardActions>
                    <TextField
                        margin="normal"
                        style={styles.messagesInput}
                        value={this.state.value}
                        onChange={this.handleChange}
                        onKeyPress={this.catchReturn}
                    /> 
                    <Button variant="raised" color="default" >
                        <input style={styles.upload} type="file" ref={this.setFileName} onChange={this.handleUploadChange}/>
                    </Button>
                    <Button size="large" style={styles.sendMessageButton} onClick={this.handleNewMessage} color="primary" disabled={disabled}>
                        Send
                    </Button>
                </CardActions>
            </Card>
        )
    }
}

export default withCookies(MessageInputAllow);