import React from 'react';
import Grid from 'material-ui/Grid';
import Paper from 'material-ui/Paper';
import Typography from 'material-ui/Typography';
import Tweet from './Tweet';

const styles = {
    logo: {
		height: 150,
	},
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
                <img style={styles.logo} alt="Logo" src={require('../Assets/Images/Twitter_Logo_Blue.png')} /> 
                    {/* <Typography variant="headline" component="h3">
                        Tweet NOW!
                    </Typography> */}
                </Paper>
            </Grid>
        );
    }
    else{
        const tweetObjs = tweetList.map((tweet) =>
            <Grid key={tweet.tweet.id} item>
                <Tweet key={tweet.tweet.id} tweet={tweet}/>
            </Grid>
        );
        return (
            tweetObjs
        );
    }
  }

  export default TweetList;