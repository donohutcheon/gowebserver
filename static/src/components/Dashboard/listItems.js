import React from 'react';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import DashboardIcon from '@material-ui/icons/Dashboard';
import SettingsIcon from '@material-ui/icons/Settings';
import PropTypes from "prop-types";

export default function MainMenuList(props) {
    const {classes, onItemClicked} = props;

    const onClick = (e, v) => {
        onItemClicked(e)
    }

    return (
        <React.Fragment>
            <ListItem id={'card-transactions-menu-item'}
                      button
                      onClick={((e) => onClick(e))}>
                <ListItemIcon>
                    <DashboardIcon/>
                </ListItemIcon>
                <ListItemText primary="Card Transactions"/>
            </ListItem>
            <ListItem id={'settings-menu-item'}
                      button
                      onClick={((e) => onClick(e))}>
                <ListItemIcon>
                    <SettingsIcon/>
                </ListItemIcon>
                <ListItemText primary="Settings"/>
            </ListItem>
        </React.Fragment>
    )
};

MainMenuList.propTypes = {
    classes: PropTypes.object.isRequired,
    onItemClicked: PropTypes.func.isRequired,
}