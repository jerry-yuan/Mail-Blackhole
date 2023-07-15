import AppHeader from "./AppHeader";
import AppContent from "./AppContent";
import { AppFooter } from "./AppFooter";

export default function AppLayout() {
    return (
        <div className="app-layout">
            <AppHeader />
            <AppContent />
            <AppFooter />
        </div>
    );
}
