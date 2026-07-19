import type { RouteObject } from "react-router";

import { AdminLayout } from "./components/admin-layout";
import { RootErrorBoundary } from "./root";
import Events from "./routes/events";
import Users from "./routes/users";
import Videos from "./routes/videos";

export const routes: RouteObject[] = [
  {
    element: <AdminLayout />,
    errorElement: <RootErrorBoundary />,
    children: [
      { index: true, element: <Videos /> },
      { path: "videos", element: <Videos /> },
      { path: "users", element: <Users /> },
      { path: "events", element: <Events /> },
    ],
  },
];
