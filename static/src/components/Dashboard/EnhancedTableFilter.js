import {makeStyles} from "@material-ui/core/styles";
import React, {useState} from "react";
import ExpansionPanel from "@material-ui/core/ExpansionPanel";
import ExpansionPanelSummary from "@material-ui/core/ExpansionPanelSummary";
import Typography from "@material-ui/core/Typography";
import ExpansionPanelDetails from "@material-ui/core/ExpansionPanelDetails";
import Autocomplete from "@material-ui/lab/Autocomplete";
import TextField from "@material-ui/core/TextField";
import Divider from "@material-ui/core/Divider";
import ExpansionPanelActions from "@material-ui/core/ExpansionPanelActions";
import Button from "@material-ui/core/Button";
import PropTypes from "prop-types";
import RangeSlider from "./RangeSlider";

const filterHeight = 240;
const useFilterStyles = makeStyles((theme) => ({
	root: {
		width: '100%',
	},
	heading: {
		fontSize: theme.typography.pxToRem(15),
		color: theme.palette.text.secondary,
	},
	details: {
		alignItems: 'center',
	},
	filter: {
		position: 'relative',
		whiteSpace: 'nowrap',
		height: filterHeight,
		transition: theme.transitions.create('height', {
			easing: theme.transitions.easing.sharp,
			duration: theme.transitions.duration.enteringScreen,
		}),
	},
	filterClose: {
		overflowY: 'hidden',
		transition: theme.transitions.create('height', {
			easing: theme.transitions.easing.sharp,
			duration: theme.transitions.duration.leavingScreen,
		}),
		height: theme.spacing(7),
		[theme.breakpoints.up('sm')]: {
			height: theme.spacing(9),
		},
	},
}));

function formatCurrency(value) {
	return (Math.round(value) / 100).toFixed(2).toLocaleString();
};

export default function EnhancedTableFilter(props) {
	const {visible} = props;
	const classes = useFilterStyles();
	const [reference, setReference] = useState("")
	const [amount, setAmount] = useState([2000, 30000])
	const formatAmount = (value) => {
		return value
	};
	const formatAmountLabel = () => {
		const label = "R" + formatCurrency(amount[0]) + " - R" + formatCurrency(amount[1])
		return label
	};
	const handleAmountChange = (event, newValue) => {
		//setAmount(newValue);
	};
	return (
		<React.Fragment>
			<ExpansionPanel
				defaultExpanded={false}
				expanded={visible}
			>
				<ExpansionPanelSummary
					aria-controls="panel1c-content"
					id="panel1c-header"
				>
					<div className={classes.column}>
						<Typography className={classes.heading}>Filter Criteria</Typography>
					</div>
				</ExpansionPanelSummary>
				<ExpansionPanelDetails className={classes.details}>
					<div className={classes.column}>
						<RangeSlider
							id="amount"
							label={"Amount"}
							value={amount}
							setValue={setAmount}
							max={100000000}
							min={100}
							onChange={handleAmountChange}
							formatAmountLabel={formatAmountLabel}
							formatValue={formatAmount}
							formatLabel={formatAmountLabel}
						/>
					</div>
					<div className={classes.column}>
						<Autocomplete
							id="reference"
							fullWidth
							options={['production','simulation']}
							getOptionLabel={(option) => option}
							style={{ width: 300 }}
							renderInput={(params) => <TextField {...params} label="Reference" variant="outlined" />}
							onChange={(e) => setReference(e.target.value)}
						/>
					</div>
				</ExpansionPanelDetails>
				<Divider />
				<ExpansionPanelActions>
					<Button size="small">Clear</Button>
					<Button size="small" color="primary">Apply</Button>
				</ExpansionPanelActions>
			</ExpansionPanel>
		</React.Fragment>
	);
};

EnhancedTableFilter.propTypes = {
	visible: PropTypes.bool.isRequired,
};