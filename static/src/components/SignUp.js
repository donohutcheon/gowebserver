import React, {useEffect, useState} from 'react';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import FormLabel from '@material-ui/core/FormLabel';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import Alert from '@material-ui/lab/Alert';
import Link from '@material-ui/core/Link';
import Grid from '@material-ui/core/Grid';
import Box from '@material-ui/core/Box';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import FormHelperText from '@material-ui/core/FormHelperText';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import useApiRequest from "../hooks/useApiRequest/useApiRequest";
import {SUCCESS, ERROR} from "../hooks/useApiRequest/actionTypes";

function Copyright() {
	return (
		<Typography variant="body2" color="textSecondary" align="center">
			{'Copyright © '}
			<Link color="inherit" href="https://material-ui.com/">
				Your Website
			</Link>{' '}
			{new Date().getFullYear()}
			{'.'}
		</Typography>
	);
}

const useStyles = makeStyles((theme) => ({
	paper: {
		marginTop: theme.spacing(8),
		display: 'flex',
		flexDirection: 'column',
		alignItems: 'center',
	},
	avatar: {
		margin: theme.spacing(1),
		backgroundColor: theme.palette.secondary.main,
	},
	form: {
		width: '100%', // Fix IE 11 issue.
		marginTop: theme.spacing(3),
	},
	submit: {
		margin: theme.spacing(3, 0, 2),
	},
}));

export default function SignUp() {
	const classes = useStyles()
	const [firstName, setFirstName] = useState('')
	const [surname, setSurname] = useState('')
	const [email, setEmail] = useState('')
	const [password, setPassword] = useState('')
	const [responseMessage, setResponseMessage] = useState('')
	const [errorFields, setErrorFields] = useState([])
	const [{ status, response }, makeRequest] = useApiRequest(
		`/api/auth/sign-up`,
		{
			verb: "post",
			data: {email, password, firstName, surname}
		}
	);
	useEffect(() => {
		if(status === SUCCESS) {
			console.log("sign-up use effect", status, response.data)
			setResponseMessage(response.data.message)
		} else if(status === ERROR) {
			console.log("sign-up use effect error", status, response.data)
			setResponseMessage(response.data.message)
			setErrorFields(response.data.fields)
		}
	}, [status, response])
	const handleSubmit = (e) => {
		e.preventDefault()
		makeRequest();
	}
	const getFieldError = (name) => {
		for (const field of errorFields) {
			if(name === field.name) {
				console.log("field is in error", name, field.message)
				return field.message
			}
		}
		return null;
	}
	return (
		<Container component="main" maxWidth="xs">
			<CssBaseline />
			<div className={classes.paper}>
				<Avatar className={classes.avatar}>
					<LockOutlinedIcon />
				</Avatar>
				<Typography component="h1" variant="h5">
					Sign up
				</Typography>
				<form className={classes.form} noValidate
				      onSubmit={handleSubmit}>
					<Grid container spacing={2}>
						<Grid item xs={12} sm={6}>
							<TextField
								error={getFieldError("firstname")?true:false}
								autoComplete="firstname"
								name="firstName"
								variant="outlined"
								required
								fullWidth
								id="firstName"
								label="First Name"
								autoFocus
								onChange={(e) => setFirstName(e.target.value)}
							/>
						</Grid>
						<Grid item xs={12} sm={6}>
							<TextField
								error={getFieldError("surname")?true:false}
								variant="outlined"
								required
								fullWidth
								id="surname"
								label="Surname"
								name="surname"
								autoComplete="surname"
								onChange={(e) => setSurname(e.target.value)}
							/>
						</Grid>
						<Grid item xs={12}>
							<TextField
								error={getFieldError("email")?true:false}
								variant="outlined"
								required
								fullWidth
								id="email"
								label="Email Address"
								name="email"
								autoComplete="email"
								onChange={(e) => setEmail(e.target.value)}
							/>
							{getFieldError("email") &&
							<FormHelperText error={true}>
								{getFieldError("email")}
							</FormHelperText>
							}
						</Grid>

						<Grid item xs={12}>
							<TextField
								error={getFieldError("password")?true:false}
								variant="outlined"
								required
								fullWidth
								name="password"
								label="Password"
								type="password"
								id="password"
								autoComplete="current-password"
								onChange={(e) => setPassword(e.target.value)}
							/>
						</Grid>
						{getFieldError("password") &&
						<Grid item xs={12}>
							<FormHelperText error={true}>
								{getFieldError("password")}
							</FormHelperText>
						</Grid>
						}
					</Grid>

					<Button
						type="submit"
						fullWidth
						variant="contained"
						color="primary"
						className={classes.submit}
					>
						Sign Up
					</Button>
					{responseMessage &&
					<Alert severity={errorFields.length>0?"error":"success"}>
						{responseMessage}
					</Alert>
					}

					<Grid container justify="flex-end">
						<Grid item>
							<Link href="login" variant="body2">
								Already have an account? Sign in
							</Link>
						</Grid>
					</Grid>
				</form>
			</div>
			<Box mt={5}>
				<Copyright />
			</Box>
		</Container>
	);
}