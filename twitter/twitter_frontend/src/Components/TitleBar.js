import React, {Component} from "react";
import { withCookies, Cookies } from 'react-cookie';
import { instanceOf } from 'prop-types';

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
        this.cookies = this.props.cookies;
    }

    componentWillMount() {
        this.checkLogin();
    }

    checkLogin =()=> {

        let session_id = this.cookies.get("session_id");

        if (session_id){
            this.goToHome()
        }else{

        }


    }

    goToHome =(e)=> {
		window.location = '/home';
	}

    render() {
        return (
            <a href="/index">
                <img style={styles.title_logo} alt="Logo" src={require('../Assets/Images/Twitter_Logo_Blue.png')} /> 
            </a>
        );
    }
}

export default withCookies(TitleBar);