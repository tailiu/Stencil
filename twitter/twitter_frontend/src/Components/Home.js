import React, {Component, Fragment} from "react";
import Grid from 'material-ui/Grid';
import NavBar from './NavBar';
import Tweet from './Tweet';
import UserInfo from './UserInfo';

const styles = {
    grid : {
        container : {
            marginTop: 80
        }
    },
    card: {
        card:{
            // minWidth: 400,
        },
        input:{
            width: "95%",
        },
        button: {
            width: "100%",
            backgroundColor: "#00aced",
            color: "#fff",
        }
    },
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
};

class Home extends Component {

    constructor(props) {

        super(props);

        // console.log(JSON.stringify(this.props.location.state.user))

        this.state = {
            email : '',
            name : '',
            tweet_value : '',
            value : '',
        }
        

        this.handleChange = this.handleChange.bind(this);
        this.handleChangeToTweetField = this.handleChangeToTweetField.bind(this);

    }

    handleChange(e) {
        this.setState({ value: e.target.value });
    }

    handleChangeToTweetField(e) {
        this.setState({ tweet_value: e.target.value });
    }   

    render () {
    return (
        <Fragment>
            <NavBar />
            <Grid style={styles.grid.container} container spacing={16}>
                
                <Grid item xs={1}>
                </Grid>

                <Grid item xs={3}>
                    <UserInfo user={this.props.location.state.user}/>

                </Grid>

                <Grid item xs={7}>
                    <Grid container spacing={8} direction="column" align="left">
                        <Grid item>
                            <Tweet />
                        </Grid>
                        <Grid item>
                            <Tweet />
                        </Grid>
                        <Grid item>
                            <Tweet />
                        </Grid>
                        <Grid item>
                            <Tweet />
                        </Grid>
                        <Grid item>
                            <Tweet />
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item xs={1}>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default Home;
