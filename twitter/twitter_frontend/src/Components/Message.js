import React, {Component} from "react";
import Moment from 'moment';
import Avatar from 'material-ui/Avatar';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
import Card, { CardContent, CardHeader, CardMedia } from 'material-ui/Card';

var styles = {
    photo: {
        height: "auto",
        width: '100%'
    },
    video: {
        height: "auto",
        width: '100%'
    },
    media_container: {
        textAlign: "center",
    },
    text: {
        whiteSpace: 'normal',
        wordWrap: 'break-word'
    }
}

class Message extends Component {

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    }

    constructor(props) {
        super(props);

        const { cookies } = this.props;

        this.state = {
            user_id: cookies.get('user_id'),
            base_url: "http://localhost:3000/"
        };
    }

    getLatestUpdatedDate = () => {
        return Moment(this.props.message.updated_at).format('MMMM Do, YYYY - h:mm A');
    }

    setStyle = () => {
        const message = this.props.message

        styles.cardContainer = {
            marginTop       : 5,
            marginBottom    : 20,
            float           : 'none', 
            width           : '40%',
            marginLeft      : 0,
            marginRight     : 0,
            borderRadius    : '20px'
        }

        if (message.user_id == this.state.user_id) {
            styles.cardContainer.marginLeft = 'auto'
        }

        return styles
    }

    getText = () => {
        const message = this.props.message
        var text = ''
        if (this.props.current_conversation_type == 'group') {
            text += message.name + ': '
        }
        text += this.props.message.content
        return text
    }

    getMedia = () => {
        const message = this.props.message
        if (message.message_media.url != null) {
            if (message.media_type == 'photo') {
                return (
                    <CardContent style={styles.media_container}>
                        <img style={styles.photo} src={this.state.base_url + this.props.message.message_media.url} />
                    </CardContent>
                )
            } else if (message.media_type == 'video') {
                return (
                    <CardContent style={styles.media_container}>
                        <video style={styles.video} controls>
                            <source src={this.state.base_url + this.props.message.message_media.url} type="video/mp4"/>
                        </video>
                    </CardContent>
                )
            }
        }
    }

    getAvatar = () => {
        return this.props.message.name.charAt(0).toUpperCase()
    }

    render () {
        const styles = this.setStyle()

        return (
            <Card style={styles.cardContainer}>
                <CardHeader
                    avatar={
                        <Avatar>
                            {this.getAvatar()}
                        </Avatar>
                    }
                    subheader={this.getLatestUpdatedDate()}
                />
                {this.getMedia()}
                <CardContent style={styles.text}>
                    {this.getText()}
                </CardContent> 
            </Card>
        )
    }
}

export default withCookies(Message);