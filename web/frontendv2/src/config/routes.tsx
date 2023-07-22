import { Outlet, createBrowserRouter } from "react-router-dom";
import InboxPage from "../views/Inbox";
import AppLayout from "../layout/AppLayout";
import NotFoundPage from "../views/NotFoundPage";
import MailPage from "../views/Mail";
export default createBrowserRouter([
    {
        path: "/",
        element: (
            <AppLayout>
                <Outlet />
            </AppLayout>
        ),
        errorElement: (
            <AppLayout>
                <Outlet />
            </AppLayout>
        ),

        children: [
            {
                path: "/",
                element: <InboxPage />,
            },
            {
                path: "mail/:id",
                element: <MailPage />,
            },
            {
                path: "*",
                element: <NotFoundPage />,
            },
        ],
    },
]);
