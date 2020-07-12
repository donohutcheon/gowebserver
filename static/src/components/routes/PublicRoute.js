import React, {useContext} from 'react';
import { Route, Redirect } from 'react-router-dom';
import {AuthContext} from "../../contexts/AuthContext";

const PublicRoute = ({component: Component, restricted, ...rest}) => {
    const {isAuthenticated} = useContext(AuthContext);
    return (
        <Route {...rest} render={props => (
            isAuthenticated && restricted ?
                <Redirect to="/dashboard" />
            : <Component {...props} />
        )} />
    );
};

export default PublicRoute;