import { Outlet } from "react-router-dom";
export default function AppContent() {
	return (
		<div className="app-content">
			<Outlet />
		</div>
	);
}
