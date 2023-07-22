import AppHeader from "./AppHeader";
import { AppFooter } from "./AppFooter";
import { ErrorBoundary } from "./ErrorBoundary";

interface Props {
    children: React.ReactNode[] | React.ReactNode;
}

const AppLayout: React.FC<Props> = function ({ children }) {
    return (
        <div className="app-layout">
            <AppHeader />
            <div className="app-content">
                <ErrorBoundary>{children}</ErrorBoundary>
            </div>
            <AppFooter />
        </div>
    );
};

export default AppLayout;
