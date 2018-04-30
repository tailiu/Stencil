import React, {Component, Fragment} from "react";
import Grid from 'material-ui/Grid';
import NavBar from './NavBar';
import MessageBar from './MessageBar';
import Tweet from './Tweet';
import UserInfo from './UserInfo';
import axios from 'axios';
import { instanceOf } from 'prop-types';
import { withCookies, Cookies } from 'react-cookie';

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

    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
      };

    constructor(props) {

        super(props);
        const { cookies } = this.props;

        this.state = {
            user_id: cookies.get('user_id'),
            email : '',
            name : '',
            tweet_value : '',
            value : '',
            user: ''
        } 
    }

    componentDidMount(){
        
    }

    goToIndex = () => {
        this.props.history.push({pathname: '/'});
      }

    handleChange =(e)=> {
        this.setState({ value: e.target.value });
    }

    handleChangeToTweetField =(e)=> {
        this.setState({ tweet_value: e.target.value });
    }   

    render () {
    return (
        <Fragment>
            <NavBar />
            <MessageBar ref={instance => { this.MessageBar = instance; }}/>
            <Grid style={styles.grid.container} container spacing={16}>
                
                <Grid item xs={1}>
                </Grid>

                <Grid item xs={3}>
                    <UserInfo user={this.state.user}/>

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

export default withCookies(Home);
