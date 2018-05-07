import React, {Component, Fragment} from 'react';
import Card, { CardContent } from "material-ui/Card";
import axios from 'axios';
import Divider from 'material-ui/Divider';
import Typography from 'material-ui/Typography';

const styles = {
    searchbox: {
        maxHeight: 300,
        marginTop: -20,
        position: "fixed",
        zIndex: 10,
        overflow: "scroll",
        right: 0,
        marginRight: 80,
        maxWidth: 300,
        minWidth: 300,
        backgroundColor: "#00aced"
    },
    user_item: {
        padding:10,
        textDecoration: "none",
        color: "#fff",
    },
    logo: {
        height: 100
    },
}

function UserList(props) {
    const userList = props.users;

    if (userList.length <= 0){
        return (
            <img style={styles.logo} alt="Logo" src={require('../Assets/Images/Twitter_Logo_Blue.png')} /> 
        );
    }
    else{
        const userObjs = userList.map((user) =>
            <a 
                style={styles.user_item}
                href={"/profile/"+user.id}>
                <Typography variant="button">
                    {user.name}, @{user.handle}
                </Typography>
            </a>
        );
        return (
            userObjs
        );
    }
  }

class UserSearchBox extends Component {

    constructor(props) {
        super(props);
        this.state = {
            search_users: [],
            query: ""
        }
    }

    componentWillReceiveProps(newProps) {
        this.setState(
            {query: newProps.query},
            () => {
                if (this.state.query === "")
                    this.setState({
                        search_users : []
                    })
                else
                    this.searchUser()
            }
        );
        console.log("HERE I AM! PROPS!")
        console.log(this.state)
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        // this.setState({query: prevProps.query});
        // console.log("HERE I AM! UPDAYED!")
        // console.log("HERE I AM! PROPS!")
        // console.log("prevProps:"+JSON.stringify(prevProps))
        // console.log("prevState:"+JSON.stringify(prevState))
        // console.log("snapshot:"+JSON.stringify(snapshot))
    }

    searchUser =()=> {
        // console.log(this.props.query)

        axios.get(
            'http://localhost:3000/users/search',
            {
                params: {
                'query': this.state.query 
                }
            }
            ).then(response => {
    
                if(response.data.result.success){
                    console.log("API:"+JSON.stringify(response.data.result.users))
                    this.setState({
                        search_users: response.data.result.users,
                    })
                }else{
                }
            })
    }

    render() {
        return (
            <Fragment>
            {this.state.search_users.length > 0 &&
                <Card style={styles.searchbox}>
                    <CardContent>
                        <UserList users={this.state.search_users} />
                    </CardContent>
                </Card>
            }
            </Fragment>
        );
    }
}

export default UserSearchBox;