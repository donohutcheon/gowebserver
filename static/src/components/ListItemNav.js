import { Link } from 'react-router-dom';

function ListItemNav(props) {
	const { icon, primary, to } = props;

	const CustomLink = props => <Link to={to} {...props} />;

	return (
		<li>
			<ListItem button component={CustomLink}>
				<ListItemIcon>{icon}</ListItemIcon>
				<ListItemText primary={primary} />
			</ListItem>
		</li>
	);
}