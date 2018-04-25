import React, {Component, Fragment} from "react";
import Grid from 'material-ui/Grid';
import NavBar from './NavBar';
import Tweet from './Tweet';

const styles = {
    grid : {
        container : {
            marginTop: 80
        }
    },
};


class Search extends Component {

    constructor(props) {

        super(props);
        this.state = {
        }
    }

  render () {
    return (
        <Fragment>
            <NavBar />
            <Grid style={styles.grid.container} container spacing={16}>
                
                <Grid item xs={2}>
                </Grid>

                <Grid item xs={8}>
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

                <Grid item xs={2}>
                </Grid>
            </Grid>
        </Fragment>
    );
  }
}

export default Search;
