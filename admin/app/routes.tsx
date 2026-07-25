import { Outlet, type RouteObject } from "react-router";

import { RootErrorBoundary } from "./root";
import EventDetail from "./routes/events/event-detail";
import Events from "./routes/events/events";
import Users from "./routes/users";
import Videos from "./routes/videos";
import { EditEvent } from "~/routes/events/edit-event";

export const routes: RouteObject[] = [
  {
    element: <Outlet />,
    errorElement: <RootErrorBoundary />,
    children: [
      { index: true, element: <Videos /> },
      { path: "videos", element: <Videos /> },
      { path: "users", element: <Users /> },
      { path: "events", element: <Events /> },
      { path: "events/:id", element: <EventDetail /> },
      { path: "events/:id/edit", element: <EditEvent /> },
    ],
  },
];
