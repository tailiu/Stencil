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
    
    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    };

    constructor(props) {
        super(props);
    }

    componentWillMount() {
        this.checkLogin();
    }

    checkLogin =()=> {
        const { cookies } = this.props;
        let session_id = cookies.get("session_id");
        console.log("session:"+ session_id)
        if (session_id){
            console.log("logged in")
            this.goToHome()
        }else{
            console.log("not logged in")
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