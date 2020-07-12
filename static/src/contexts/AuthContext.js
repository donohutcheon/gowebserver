import React, {createContext, useEffect, useState} from 'react';
export const AuthContext = createContext();

const AuthContextProvider = (props) => {
	const parseJwt = (token) => {
		const base64Url = token.split('.')[1];
		const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
		const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
			return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
		}).join(''));

		return JSON.parse(jsonPayload);
	};
	const [rememberMe, setRememberMe] = useState(false);
	const [accessToken, setAccessToken] = useState(() => {
		const localData = localStorage.getItem('accessToken')
		return localData ? localData : "";
	});

	const [refreshToken, setRefreshToken] = useState(() => {
		const localData = localStorage.getItem('refreshToken')
		return localData ? localData : "";
	});

	const [isAuthenticated, setAuthenticated] = useState(() => {
		let localData = rememberMe ? localStorage.getItem('accessToken') : accessToken
		if(!localData) {
			return false
		}
		const tokenExpiry = parseJwt(localData).exp
		console.log(new Date().getTime(), tokenExpiry * 1000)
		console.log(new Date().getTime() <  tokenExpiry * 1000 ? "Authenticated" : "Not Authenticated")
		return new Date().getTime() <  tokenExpiry * 1000
	});

	const setAuthentication = (access, refresh, rememberLogin) => {
		console.log(access)
		setRememberMe(rememberLogin);
		setAccessToken(access)
		setRefreshToken(refresh)
		const tokenExpiry = parseJwt(access).exp
		setAuthenticated(new Date().getTime() <  tokenExpiry * 1000)
		if(rememberLogin) {
			localStorage.setItem('accessToken', access)
			localStorage.setItem('refreshToken', refresh)
		}
	}

	/*useEffect(() => {
		console.log(accessToken)
		const exp = parseJwt(accessToken).exp
		setAuthenticated(new Date().getTime() <  exp * 1000)
		localStorage.setItem('accessToken', accessToken)
	}, [accessToken])*/

	/*useEffect(() => {
		localStorage.setItem('refreshToken', refreshToken)
	}, [refreshToken])*/

	return (
		<AuthContext.Provider value={{isAuthenticated, accessToken, refreshToken, setAuthentication}}>
			{props.children}
		</AuthContext.Provider>
	);
}

export default AuthContextProvider;