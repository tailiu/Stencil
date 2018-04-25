import React, {Component, Fragment} from "react";
import Grid from 'material-ui/Grid';
import Card, { CardContent, CardHeader } from 'material-ui/Card';
import NavBar from './NavBar';
import List, { ListItem, ListItemIcon, ListItemText } from 'material-ui/List';

const styles = {
    grid : {
        container : {
            marginTop: 80,
            height: 200
        }
    },
};

class Notif extends Component {

    constructor(props) {

        super(props);
        this.state = {
        }
    }

  render () {
    return (
        <Fragment>
            <NavBar />
            <Grid style={styles.grid.container} direction="column" container spacing={24} >
                
                <Grid item xs={2}>
                </Grid>
                <Grid item xs={8}>


                            <Card>
                                <CardHeader
                                    title="Notifications"
                                />
                                <hr />
                                <CardContent>

                                    <List dense={true}>
                                        <ListItem>
                                            <ListItemIcon>
                                            <img style={styles.logo} alt="retweet" src={require('../Assets/Images/retweet_icon.png')} /> 
                                            </ListItemIcon>
                                            <ListItemText
                                            primary="Tai retweeted your tweet"
                                            //   secondary={secondary ? 'Secondary text' : null}
                                            />
                                        </ListItem>
                                        <ListItem>
                                            <ListItemIcon>
                                            <img style={styles.logo} alt="favourite" src={require('../Assets/Images/fav_icon.png')} /> 
                                            </ListItemIcon>
                                            <ListItemText
                                            primary="Miro favorited your tweet"
                                            //   secondary={secondary ? 'Secondary text' : null}
                                            />
                                        </ListItem>
                                        <ListItem>
                                            <ListItemIcon>
                                            <img style={styles.logo} alt="follow" src={require('../Assets/Images/follow_icon.png')} /> 
                                            {/* <FolderIcon /> */}
                                            </ListItemIcon>
                                            <ListItemText
                                            primary="Major Tom followed you"
                                            //   secondary={secondary ? 'Secondary text' : null}
                                            />
                                        </ListItem>
                                    </List>
                                </CardContent>
                            </Card>


                </Grid>
                <Grid item xs={2}>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default Notif;
