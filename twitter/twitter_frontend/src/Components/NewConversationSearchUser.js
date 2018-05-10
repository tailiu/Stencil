import React from "react";
import Downshift from 'downshift';
import Paper from 'material-ui/Paper';
import { MenuItem } from 'material-ui/Menu';
import Chip from 'material-ui/Chip';
import TextField from 'material-ui/TextField';
import keycode from 'keycode';
import { withStyles } from 'material-ui/styles';
import PropTypes from 'prop-types';

const styles = theme => ({
    root: {
        flexGrow: 1,
        height: 250
    },
    container: {
        flexGrow: 1,
        position: 'relative'
    },
    paper: {
        position: 'absolute',
        zIndex: 1,
        marginTop: theme.spacing.unit,
        left: 0,
        right: 0
    },
    chip: {
        margin: `${theme.spacing.unit / 2}px ${theme.spacing.unit / 4}px`
    },
    inputRoot: {
        flexWrap: 'wrap'
    },
});

class NewConversationSearchUser extends React.Component {

    constructor(props) {

        super(props);

        this.state = {
            inputValue: '',
            selectedItem: [],
        };
    }

    renderInput = inputProps =>  {
        const { InputProps, classes, ref, ...other } = inputProps;
    
        return (
            <TextField
                InputProps={{
                inputRef: ref,
                classes: {
                    root: classes.inputRoot,
                },
                ...InputProps,
                }}
                {...other}
            />
        );
    }
    
    renderSuggestion = ({ suggestion, index, itemProps, highlightedIndex, selectedItem }) => {
        const isHighlighted = highlightedIndex === index;
        const isSelected = (selectedItem || '').indexOf(suggestion) > -1;
    
        return (
            <MenuItem
                {...itemProps}
                key={suggestion}
                selected={isHighlighted}
                component="div"
                style={{
                fontWeight: isSelected ? 500 : 400,
                }}
            >
                {suggestion}
            </MenuItem>
        );
    }

    getSuggestions = (inputValue) => {
        let count = 0;
        
        const suggestions = this.props.suggestions

        return suggestions.filter(suggestion => {
            const keep =
                (!inputValue || suggestion.toLowerCase().indexOf(inputValue.toLowerCase()) !== -1) &&
                count < 5;
    
            if (keep) {
                count += 1;
            }
    
            return keep;
        });
    }

    handleKeyDown = event => {
        const { inputValue, selectedItem } = this.state;
        if (selectedItem.length && !inputValue.length && keycode(event) === 'backspace') {
            this.setState({
            selectedItem: selectedItem.slice(0, selectedItem.length - 1),
            });
        }
    };

    handleInputChange = event => {
        this.setState({ inputValue: event.target.value });
    };

    handleChange = item => {
        let { selectedItem } = this.state;

        if (selectedItem.indexOf(item) === -1) {
            selectedItem = [...selectedItem, item];
        }

        this.setState({
            inputValue: '',
            selectedItem,
        });
    };

    handleDelete = item => () => {
        const selectedItem = [...this.state.selectedItem];
        selectedItem.splice(selectedItem.indexOf(item), 1);

        this.setState({ selectedItem });
    };

    render() {
        const { classes } = this.props;
        const { inputValue, selectedItem } = this.state;
        console.log(classes)

        return (
            <Downshift inputValue={inputValue} onChange={this.handleChange} selectedItem={selectedItem}>
                {({
                    getInputProps,
                    getItemProps,
                    isOpen,
                    inputValue: inputValue2,
                    selectedItem: selectedItem2,
                    highlightedIndex,
                }) => (
                    <div className={classes.container}>
                        {this.renderInput
                            ({
                                fullWidth: true,
                                classes,
                                InputProps: getInputProps({
                                startAdornment: selectedItem.map(item => (
                                    <Chip
                                        key={item}
                                        tabIndex={-1}
                                        label={item}
                                        onDelete={this.handleDelete(item)}
                                    />
                                )),
                                onChange: this.handleInputChange,
                                onKeyDown: this.handleKeyDown,
                                placeholder: 'Enter one or multiple names',
                                id: 'integration-downshift-multiple',
                            }),
                        })}
                        {isOpen ? 
                        (
                            <Paper className={classes.paper} square>
                                {this.getSuggestions(inputValue).map((suggestion, index) =>
                                    this.renderSuggestion({
                                        suggestion,
                                        index,
                                        itemProps: getItemProps({ item: suggestion }),
                                        highlightedIndex,
                                        selectedItem,
                                    }),
                                )}
                            </Paper>
                        ) : null}
                    </div>
                )}
            </Downshift>
        );
    }
}

NewConversationSearchUser.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(NewConversationSearchUser);
