import React, {Component, Fragment} from "react";
import Button from 'material-ui/Button';
import Avatar from 'material-ui/Avatar';
import Typography from 'material-ui/Typography';
import TextField from 'material-ui/TextField';
import Card, { CardActions, CardContent, CardHeader } from 'material-ui/Card';
import Dialog, {
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
  } from 'material-ui/Dialog';

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

class UserProfileBox extends Component{

    constructor(props){
        super(props);
        this.state = {
            bio_box_open : false
        }
    }

    handleBioBoxOpen = () => {
        console.log("HERE!");
        this.setState({bio_box_open: true });
    };

    handleBioBoxClose = () => {
        this.setState({ bio_box_open: false });
    };

    render(){
        return(
            <Fragment>
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
                <CardContent>
                    <Typography>
                        Your bio here!
                    </Typography>
                </CardContent>
                <CardActions>
                    <Button size="small" onClick={this.like}>
                        Change Photo
                    </Button>
                    <Button size="small" onClick={this.handleBioBoxOpen}>
                        Change Bio
                    </Button>
                </CardActions>
            </Card>
            <Dialog
                open={this.state.bio_box_open}
                onClose={this.handleBioBoxOpen}
                aria-labelledby="form-dialog-title"
                >
                <DialogTitle id="form-dialog-title">Change Bio</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                    {/* What's on your mind? */}
                    </DialogContentText>
                    <TextField
                    autoFocus
                    id="bio"
                    name="bio"
                    // label="What's on your mind?"
                    fullWidth
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={this.handleBioBoxClose} color="primary">
                    Cancel
                    </Button>
                    <Button onClick={this.handleBioBoxClose} color="primary">
                    Change
                    </Button>
                </DialogActions>
            </Dialog>
            </Fragment>
        );
    }
}

export default UserProfileBox;