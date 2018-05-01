import React from 'react';
import Grid from 'material-ui/Grid';
import Paper from 'material-ui/Paper';
import Typography from 'material-ui/Typography';
import Tweet from './Tweet';

const styles = {
    paper: {
        padding: 15,
        textAlign : "center",
        minHeight: 49
    }
}

function TweetList(props) {
    const tweetList = props.tweets;

    if (tweetList.length <= 0){
        return (
            <Grid item>
                <Paper elevation={2} style={styles.paper}>
                    <Typography variant="headline" component="h3">
                        No Tweets Yet
                    </Typography>
                </Paper>
            </Grid>
        );
    }
    else{
        const tweetObjs = tweetList.map((tweet) =>
            <Grid key={tweet.id} item>
                <Tweet key={tweet.id} tweet={tweet}/>
            </Grid>
        );
        return (
            tweetObjs
        );
    }
  }

  export default TweetList;