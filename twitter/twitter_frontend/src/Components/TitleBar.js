import React, {Component} from "react";

// import Typography from 'material-ui/Typography';

// import AppBar from 'material-ui/AppBar';
// import Toolbar from 'material-ui/Toolbar';

const styles = {
    titlebar: {
        backgroundColor: "#00aced",
    },
    title: {
        color: "#fff",
        cursor: "pointer",
    },
    title_logo: {
        cursor: "pointer",
        height: 150,
    }
};

class TitleBar extends Component {

    constructor(props) {
        super(props);
        this.goToHome = this.goToHome.bind(this);
    }


    goToHome(e) {
		window.location = '/index';
	}

    render() {
        return (
            <a href="/index">
                <img style={styles.logo} alt="Logo" src={require('../Assets/Images/Twitter_Logo_Blue.png')} /> 
            </a>
            // <AppBar 
            //     style={styles.titlebar} 
            //     position="static" 
            //     color="default">
            //     <Toolbar>
            //         <Typography 
            //             variant="title" 
            //             style={styles.title} 
            //             onClick = {this.goToHome}>
            //             Twitter
            //         </Typography>
            //     </Toolbar>
            // </AppBar>
        );
    }
}

export default TitleBar;