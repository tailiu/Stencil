import React, {Component} from "react";
import Button from 'material-ui/Button';
import Avatar from 'material-ui/Avatar';
import Typography from 'material-ui/Typography';
import Card, { CardActions, CardContent, CardHeader } from 'material-ui/Card';
import Moment from 'moment';
import renderHTML from 'react-render-html';
import IconButton from 'material-ui/IconButton';
import axios from 'axios';
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';
import NewTweetDialog from './NewTweetDialog';

const styles = {
    tweet: {
        avatar: {

        },
        main_input: {

        },
        container: {

        },
        card: {

        },
        actions: {
            flex: 1
        },
        goto_icon:{
            height:15,
            opacity:0.7
        },
        action_icon: {
            height:22,
            // opacity:0.7
        },
        action_stat: {
            display: "inline-block",
            // color:"red",
            opacity: 0.7,
            fontSize: 15,
            fontFamily: '"Courier New", Courier, "Lucida Sans Typewriter"'
        },
        image: {
            height: 240,
            textAlign: "center"
        },
        image_area: {
            height: 240,
            alignItems: "center",
            textAlign: "center",
            marginTop: 20,
            background: "#e6ecf0"
        }
    }
}


class Tweet extends Component{

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    };
    
    constructor(props){
        super(props);

        const { cookies } = this.props;
        // console.log(props.tweet.tweet.id);
        this.state = {
            base_url : "http://localhost:3000/",
            user_id: parseInt(cookies.get('user_id')),
            user_name: cookies.get('user_name'),
            user_handle: cookies.get('user_handle'),
            liked: false,
            retweeted: false,
            replied: false,
            likes:this.props.tweet.likes,
            retweets:this.props.tweet.retweets,
            replies:this.props.tweet.replies,
            tweet_box_open: false,

        }
    }

    componentDidMount(){
        this.stats();
    }

    handleTweetBoxClose = () => {
        this.setState({ tweet_box_open: false });
        this.stats();
    };

    handleTweetBoxOpen = () => {
        console.log("HERE!");
        this.setState({tweet_box_open: true });
    };


    stats = () => {
        axios.get(
        'http://localhost:3000/tweets/stats',
        {
            params: {
            'tweet_id': this.props.tweet.tweet.id, 
            }
        }
        ).then(response => {
        if(response.data.result.success){
            this.setState({
                likes: response.data.result.likes.length,
                retweets: response.data.result.retweets.length,
                replies: response.data.result.replies.length,
            })

            if(response.data.result.likes.indexOf(this.state.user_id)>=0){
                this.setState({
                    liked: true
                })
            }

            if(response.data.result.retweets.indexOf(this.state.user_id)>=0){
                this.setState({
                    retweeted: true
                })
            }
            
            if(response.data.result.replies.indexOf(this.state.user_id)>=0){
                this.setState({
                    replied: true
                })
            }

        }else{
            console.log("Unable to fetch stats!")
        }
        })
    }

    like = (like,e) => {
        axios.get(
        'http://localhost:3000/tweets/like',
        {
            params: {
            'user_id': this.state.user_id, 
            'tweet_id': this.props.tweet.tweet.id, 
            'like': like
            }
        }
        ).then(response => {
        if(response.data.result.success){
            this.setState({
                liked: like
            })
            this.stats();
        }else{
            console.log("Unable to like!")
        }
        })
    }

    retweet = (retweet, e) => {
        axios.get(
            'http://localhost:3000/tweets/retweet',
            {
                params: {
                'user_id': this.state.user_id, 
                'tweet_id': this.props.tweet.tweet.id, 
                'retweet': retweet
                }
            }
            ).then(response => {
            if(response.data.result.success){
                this.setState({
                    retweeted: retweet
                })
                this.stats();
            }else{
                console.log("Unable to retweet!")
            }
            })
    }

    delete = (e) => {
        axios.get(
            'http://localhost:3000/tweets/delete',
            {
                params: {
                'user_id': this.state.user_id, 
                'tweet_id': this.props.tweet.tweet.id, 
                }
            }
            ).then(response => {
            if(response.data.result.success){
                window.location.reload();
            }else{
                console.log("Unable to retweet!")
            }
            })
    }

    reply = e => {
        
    }

    render(){
        Moment.locale('en');

        return(
            <Card id={this.props.tweet.tweet.id}>
                <CardHeader
                    avatar={
                        <Avatar aria-label="Recipe" style={styles.tweet.avatar}>
                        {this.props.tweet.creator.name[0]}
                        </Avatar>
                    }
                    // onClick={this.goToProfile(this.props.tweet.creator.id)}
                    // title={this.props.tweet.creator.name}
                    title={renderHTML('<a style="text-decoration: none;" href="/profile/'+this.props.tweet.creator.id+'">'+this.props.tweet.creator.name+'</a>' )}
                    subheader={"@"+this.props.tweet.creator.handle}
                    action={
                        <IconButton>
                            <a href={"/tweet/"+this.props.tweet.tweet.id}><img style={styles.tweet.goto_icon} alt="Logo" src={require('../Assets/Images/goto-link-icon.png')} /> </a>
                        </IconButton>
                    }
                />
                <CardContent>
                    <Typography component="p">
                    {this.props.tweet.tweet.content}
                    </Typography>
                    {this.props.tweet.tweet.media_type === "photo" &&
                        <div style={styles.tweet.image_area}>
                            <img style={styles.tweet.image} src={this.state.base_url+this.props.tweet.tweet.tweet_media.url} />
                        </div>
                    }
                    {this.props.tweet.tweet.media_type === "video" &&
                        <div style={styles.tweet.image_area}>
                            <video height="240" controls>
                                <source src={this.state.base_url+this.props.tweet.tweet.tweet_media.url} type="video/mp4"/>
                                Your browser does not support the video tag.
                            </video>
                        </div>
                    }
                
                </CardContent>
                <CardActions>
                    <div style={styles.tweet.actions}>
                        {this.state.liked? 
                            <IconButton size="small" onClick={this.like.bind(this, false)}>
                                <img style={styles.tweet.action_icon} alt="Unlike" src={require('../Assets/Images/liked-icon.png')} />
                            </IconButton>
                        :
                        <IconButton size="small" onClick={this.like.bind(this, true)}>
                                <img style={styles.tweet.action_icon} alt="Like" src={require('../Assets/Images/like-icon.png')} />
                            </IconButton>
                        }
                        <div style={styles.tweet.action_stat}>
                            {this.state.likes}
                        </div>
                        {this.state.retweeted?
                            <IconButton size="small" onClick={this.retweet.bind(this, false)}>
                                <img style={styles.tweet.action_icon} alt="UnRetweet" src={require('../Assets/Images/retweeted-icon.png')} />
                            </IconButton>
                        :
                        <IconButton size="small" onClick={this.retweet.bind(this, true)}>
                                <img style={styles.tweet.action_icon} alt="Retweet" src={require('../Assets/Images/retweet-icon.png')} />
                            </IconButton>
                        }
                        <div style={styles.tweet.action_stat}>
                            {this.state.retweets}
                        </div>
                        {this.state.replied?
                            <IconButton size="small" onClick={this.handleTweetBoxOpen}>
                                <img style={styles.tweet.action_icon} alt="Reply" src={require('../Assets/Images/replied-icon.png')} />
                            </IconButton>                        
                        :
                            <IconButton size="small" onClick={this.handleTweetBoxOpen}>
                                <img style={styles.tweet.action_icon} alt="Reply" src={require('../Assets/Images/reply-icon.png')} />
                            </IconButton>
                        }
                        <div style={styles.tweet.action_stat}>
                            {this.state.replies}
                        </div>
                        {this.state.user_id === this.props.tweet.creator.id &&
                            <IconButton size="small" onClick={this.delete}>
                                <img style={styles.tweet.action_icon} alt="Reply" src={require('../Assets/Images/delete-icon.png')} />
                            </IconButton>
                        }
                    </div>
                    <Typography component="p">
                        {Moment(this.props.tweet.tweet.created_at).format('MMMM Do, YYYY - h:mm A')}
                    </Typography>
                    <NewTweetDialog open={this.state.tweet_box_open} onChange={this.handleTweetBoxClose} reply_id={this.props.tweet.tweet.id}/>
                </CardActions>
            </Card>
        );
    }
}

export default withCookies(Tweet);