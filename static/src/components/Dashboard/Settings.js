import React, {useEffect} from "react";
import Grid from "@material-ui/core/Grid";
import Label from "recharts/lib/component/Label";
import {FETCHING, SUCCESS} from "../../hooks/useApiRequest/actionTypes";
import Button from "@material-ui/core/Button";
import {makeStyles} from "@material-ui/core/styles";
import useApiRequest from "../../hooks/useApiRequest/useApiRequest";
import Divider from "@material-ui/core/Divider";
import Typography from "@material-ui/core/Typography";
import { formatISO, fromUnixTime } from 'date-fns'
import TextField from "@material-ui/core/TextField";

const useStyles = makeStyles((theme) => ({
	submit: {
		margin: theme.spacing(3, 0, 2),
	},
	title: {

	},
	value: {

	},
}));

export default function Settings() {
	const classes = useStyles();
	const [apiKeyExpiry, setAPIKeyExpiry] = React.useState("");
	const [apiToken, setAPIKey] = React.useState("");
	const [{ status, response }, makeRequest] = useApiRequest(
		`/api/auth/api-token`,
		{
			verb: "get",
		}
	);

	useEffect(() => {
		if(status === SUCCESS) {
			console.log(response.data.token.expiresIn);
			console.log(response.data.token.apiToken);
			const date = fromUnixTime(response.data.token.expiresIn)
			setAPIKeyExpiry(formatISO(date))
			setAPIKey(response.data.token.apiToken)
		}
	}, [status])

	const onFetchAPIToken = (e) => {
		makeRequest();
	}

	return (
		<React.Fragment>
			<Grid container
			      direction="column"
			      justify="flex-start"
			      alignItems="stretch"
			      spacing={4}>
				<Grid item xs={12}>
					<Grid container
					      direction="row"
					      justify="flex-start"
					      alignItems="stretch"
					      spacing={1}>
						<Grid item xs={12}>
							{apiToken.length == 0 &&
							<Typography className={classes.title} variant="h6" id="apiTokenLabel" component="div">
								Show API Token
							</Typography>
							}
							{apiToken.length > 0 &&
							<Typography className={classes.title} variant="h6" id="apiTokenLabel" component="div">
								API Token
							</Typography>
							}
						</Grid>
						{apiToken.length == 0 &&
						<Grid item xs={2}>
							<Button
								type="submit"
								fullWidth
								variant="contained"
								color="primary"
								disabled={status === FETCHING}
								onClick={onFetchAPIToken}
							>
								Show API Token
							</Button>
						</Grid>
						}
						{apiToken.length > 0 &&
						<Grid item xs={12}>
							<TextField
								variant="outlined"
								margin="normal"
								fullWidth
								InputProps={{
									readOnly: true,
								}}
								multiline="true"
								id="api-token-value"
								value={apiToken}/>
						</Grid>
						}
						{apiToken.length > 0 &&
						<Grid item xs={12}>
							<Typography className={classes.value} variant="caption" id="apiTokenCaption"
							            component="div">
								Please keep this API key secret.
							</Typography>
						</Grid>
						}
						<Divider/>
					</Grid>
				</Grid>
			</Grid>
		</React.Fragment>
	);
}