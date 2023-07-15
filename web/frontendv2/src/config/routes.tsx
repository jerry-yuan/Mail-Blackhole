import { createBrowserRouter } from "react-router-dom";
import InboxPage from "../views/Inbox";
import AppLayout from "../layout/AppLayout";
import ErrorPage from "../views/ErrorPage";
import MailPage from "../views/Mail";
export default createBrowserRouter([
    {
        path: "/",
        element: <AppLayout />,
        errorElement: <AppLayout />,

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
                element: <ErrorPage />,
            },
        ],
    },
]);
