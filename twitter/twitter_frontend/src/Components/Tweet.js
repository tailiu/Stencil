import React, {Component} from "react";
import Button from 'material-ui/Button';
import Avatar from 'material-ui/Avatar';
import Typography from 'material-ui/Typography';
import Card, { CardActions, CardContent, CardHeader } from 'material-ui/Card';
import Moment from 'moment';
import renderHTML from 'react-render-html';

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
        }
    }
}


class Tweet extends Component{

    constructor(props){
        super(props);
        console.log(props);
        this.state = {

        }
    }

    like = e => {

    }

    retweet = e => {
        
    }

    reply = e => {
        
    }

    render(){
        Moment.locale('en');
        return(
            <Card>
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
                />
                
                <CardContent>
                    <Typography component="p">
                        {this.props.tweet.tweet.content}
                    </Typography>
                
                </CardContent>
                <CardActions>
                    <div style={styles.tweet.actions}>
                        <Button size="small" onClick={this.like}>
                            Like
                        </Button>
                        <Button size="small" onClick={this.retweet}>
                            Retweet
                        </Button>
                        <Button size="small" onClick={this.reply}>
                            Reply
                        </Button>
                    </div>
                    <Typography component="p">
                        {Moment(this.props.tweet.tweet.created_at).format('MMMM Do, YYYY - h:mm A')}
                    </Typography>
                </CardActions>
            </Card>
        );
    }
}

export default Tweet;