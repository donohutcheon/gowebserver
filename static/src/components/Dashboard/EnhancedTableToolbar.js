import {lighten, makeStyles} from "@material-ui/core/styles";
import React, {useState} from "react";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import Tooltip from "@material-ui/core/Tooltip";
import IconButton from "@material-ui/core/IconButton";
import RefreshIcon from '@material-ui/icons/Refresh';
import FilterListIcon from "@material-ui/icons/FilterList";
import EnhancedTableFilter from "./EnhancedTableFilter";
import PropTypes from "prop-types";

const useToolbarStyles = makeStyles((theme) => ({
	root: {
		paddingLeft: theme.spacing(2),
		paddingRight: theme.spacing(1),
	},
	highlight:
		theme.palette.type === 'light'
			? {
				color: theme.palette.secondary.main,
				backgroundColor: lighten(theme.palette.secondary.light, 0.85),
			}
			: {
				color: theme.palette.text.primary,
				backgroundColor: theme.palette.secondary.dark,
			},
	title: {
		flex: '1 1 100%',
	},
}));

export default function EnhancedTableToolbar(props) {
	const {onRefreshTable} = props;
	const classes = useToolbarStyles();
	const [isFilterVisible, setFilterVisible] = useState(false)
	const showFilter = (property) => (event) => {
		setFilterVisible(!isFilterVisible)
	};

	return (
		<React.Fragment>
			<Toolbar>
				<Typography className={classes.title} variant="h6" id="tableTitle" component="div">
					Card Transactions
				</Typography>
				<Tooltip title="Refresh">
					<IconButton
						aria-label="refresh table"
						onClick={onRefreshTable}
					>
						<RefreshIcon/>

					</IconButton>
				</Tooltip>
				<Tooltip title="Filter list">
					<IconButton
						aria-label="filter list"
						onClick={showFilter()}
					>
						<FilterListIcon/>
					</IconButton>
				</Tooltip>
			</Toolbar>
			<EnhancedTableFilter
				visible={isFilterVisible}
			/>
		</React.Fragment>
	);
};

EnhancedTableToolbar.propTypes = {
	classes: PropTypes.object.isRequired,
	onRefreshTable: PropTypes.func.isRequired,
};
