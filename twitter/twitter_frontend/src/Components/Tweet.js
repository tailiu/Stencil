import React, {Component} from "react";
import Button from 'material-ui/Button';
import Avatar from 'material-ui/Avatar';
import Typography from 'material-ui/Typography';
import Card, { CardActions, CardContent, CardHeader } from 'material-ui/Card';

const styles = {
    tweet: {
        avatar: {

        },
        main_input: {

        },
        container: {

        },
        card: {

        }
    }
}


class Tweet extends Component{

    constructor(props){
        super(props);
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
        return(
            <Card>
                <CardHeader
                    avatar={
                        <Avatar aria-label="Recipe" style={styles.tweet.avatar}>
                        ZT
                        </Avatar>
                    }
                    title="Tai Cow"
                    subheader="September 14, 2016"
                />
                <CardContent>
                    <Typography component="p">
                        This impressive paella is a perfect party dish and a fun meal to cook together with
                    your guests. Add 1 cup of frozen peas along with the mussels, if you like.
                    </Typography>
                
                </CardContent>
                <CardActions>
                    <Button size="small" onClick={this.like}>
                        Like
                    </Button>
                    <Button size="small" onClick={this.retweet}>
                        Retweet
                    </Button>
                    <Button size="small" onClick={this.reply}>
                        Reply
                    </Button>
                </CardActions>
            </Card>
        );
    }
}

export default Tweet;