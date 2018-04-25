import React, {Component} from "react";
import Avatar from 'material-ui/Avatar';
import Typography from 'material-ui/Typography';
import Card, { CardHeader } from 'material-ui/Card';

const styles = {
    user_info: {
        avatar: {

        },
        container: {
            // backgroundColor: "#00aced",
            // color: "#fff"
        },
        card: {

        }
    }
}

class UserInfo extends Component{

    render(){
        return(
            <Card align="left" style={styles.user_info.container}>
                <CardHeader
                    avatar={
                        <Avatar aria-label="Recipe" style={styles.user_info.avatar}>
                        TC
                        </Avatar>
                    }
                    title="Tai Big Fat Cow"
                    subheader="Followers:49, Following:51, Tweets:90"
                />
            </Card>
        );
    }
}

export default UserInfo;