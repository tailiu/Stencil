import React, {Component, Fragment} from "react";
import Typography from 'material-ui/Typography';
import Card, { CardContent } from 'material-ui/Card';

const styles = {
    backgroundColor: '#E57373',
    height: '100%',
    fontSize: 20
}

class MessageInputBlock extends Component {
    constructor(props) {
        super(props)
    }
 
    render() {
        
        return (
            <Card style={styles}>
                <CardContent>
                    <Typography variant="Subheading" gutterBottom align="center">
                        You can no longer send messages to this person.
                    </Typography>
                </CardContent>
            </Card>
        )
    }
}

export default MessageInputBlock;