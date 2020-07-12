import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import React from "react";
import TableCell from "@material-ui/core/TableCell";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import PropTypes from "prop-types";

export default function EnhancedTableHead(props) {
	const {classes, order, orderBy, onRequestSort, headCells} = props;
	const createSortHandler = (property) => (event) => {
		onRequestSort(event, property);
	};

	return (
		<TableHead>
			<TableRow>
				{headCells.map((headCell) => (
					<React.Fragment>
						{headCell.canSort &&
						<TableCell
							key={headCell.id}
							align={headCell.numeric ? 'right' : 'left'}
							padding={headCell.disablePadding ? 'none' : 'default'}
							sortDirection={orderBy === headCell.id ? order : false}
						>
							<TableSortLabel
								active={orderBy === headCell.id}
								direction={orderBy === headCell.id ? order : 'asc'}
								onClick={createSortHandler(headCell.id)}
							>
								{headCell.label}
								{orderBy === headCell.id &&
								<span className={classes.visuallyHidden}>
                                {order === 'desc' ? 'sorted descending' : 'sorted ascending'}
                                </span>}
							</TableSortLabel>
						</TableCell>
						}
						{!headCell.canSort &&
						<TableCell
							key={headCell.id}
							align={headCell.numeric ? 'right' : 'left'}
							padding={headCell.disablePadding ? 'none' : 'default'}
						>
							{headCell.label}
						</TableCell>
						}
					</React.Fragment>
				))}
			</TableRow>
		</TableHead>
	);
}

EnhancedTableHead.propTypes = {
	classes: PropTypes.object.isRequired,
	onRequestSort: PropTypes.func.isRequired,
	order: PropTypes.oneOf(['asc', 'desc']).isRequired,
	orderBy: PropTypes.string.isRequired,
	headCells: PropTypes.arrayOf(PropTypes.shape({
		id: PropTypes.string.isRequired,
		numeric: PropTypes.bool.isRequired,
		disablePadding: PropTypes.bool.isRequired,
		canSort: PropTypes.bool.isRequired,
		label: PropTypes.string.isRequired,
	})).isRequired
};