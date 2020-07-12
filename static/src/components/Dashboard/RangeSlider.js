import React, {useEffect, useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Slider from '@material-ui/core/Slider';
import PropTypes from "prop-types";
import Input from "@material-ui/core/Input";
import Grid from "@material-ui/core/Grid";

const useStyles = makeStyles((theme) => ({
	root: {
		width: 300,
		padding: '20px'
	},
	input: {
		width: 42,
	},
	rangeLabel: {
		fontSize: theme.typography.pxToRem(15),
		color: theme.palette.text.secondary,
	}
}));

const power = 12;


export default function RangeSlider(props) {
	const {value, setValue, onChange, formatLabel, formatValue, label, min, max} = props;
	const classes = useStyles();
	const [rangeLabel, setRangeLabel] = useState("")
	const [proxyValue, setProxyValue] = useState([1,2])

	function transform(value) {
		return 10 ** value
	}

	function transformRange(value) {
		if (value.length < 2)
			return []
		return [transform(value[0]),transform(value[1])]
	}

	function reverse(value) {
		return Math.log10(value);
	}

	function reverseRange(value) {
		console.log("1 " + value)
		if (value.length < 2)
			return [1,2]
		console.log("2 " + value)
		return [reverse(value[0]),reverse(value[1])]
	}

	function onProxyValueChange(event, newValue) {
		setProxyValue(newValue)
	}

	useEffect(() => {
		setValue(transformRange(proxyValue))
	}, [proxyValue]);

	useEffect(() => {
		setRangeLabel(formatLabel())
	}, [value]);

	useEffect(() => {
		console.log("3")
		setProxyValue(reverseRange(value))
		console.log("4")
	}, []);

	return (
		<div className={classes.root}>
			<Typography id="range-slider-label" gutterBottom>
				{label}
			</Typography>
			<Typography id="range-slider-value-label" gutterBottom className={classes.rangeLabel}>
				{rangeLabel}
			</Typography>
			<Grid container spacing={2} alignItems="center">
				<Grid item xs>
					<Slider
						step={0.01}
						value={proxyValue}
						onChange={onProxyValueChange}
						valueLabelDisplay="off"
						aria-labelledby="range-slider-label"
						getAriaValueText={formatValue}
						min={reverse(min)}
						max={reverse(max)}
					/>
				</Grid>
			</Grid>
		</div>
	);
}

RangeSlider.propTypes = {
	value: PropTypes.arrayOf(PropTypes.number).isRequired,
	setValue: PropTypes.func.isRequired,
	formatLabel: PropTypes.func.isRequired,
	formatValue: PropTypes.func.isRequired,
	onChange: PropTypes.func.isRequired,
	label: PropTypes.string.isRequired,
	min: PropTypes.number.isRequired,
	max: PropTypes.number.isRequired
};