import React, {useContext} from 'react';
import {Link} from 'react-router-dom';
import {AuthContext} from "../contexts/AuthContext";

export default function About(props) {
	const {isAuthenticated} = useContext(AuthContext)
	//console.log(isAuthenticated)
	return (
		<div>
			<h1>About Page</h1>
			<p>{isAuthenticated ? "Authenticated" : "Not Authenticated"}</p>
			<Link to="/login">Login</Link>
			<br/>
			<Link to="/dashboard">Dashboard</Link>
		</div>
	);
}