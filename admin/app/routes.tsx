import { Outlet, type RouteObject } from "react-router";

import { RootErrorBoundary } from "./root";
import EventDetail from "./routes/events/event-detail";
import Events from "./routes/events/events";
import NewEvent from "./routes/events/new-event";
import Users from "./routes/users";
import VideoGroups from "./routes/video-groups/video-groups";
import NewVideoGroup from "./routes/video-groups/new-video-group";
import { EditVideoGroup } from "./routes/video-groups/edit-video-group";
import { MarkdownEditorDialog } from "~/components/markdown-editor";
import { EditEvent } from "~/routes/events/edit-event";

export const routes: RouteObject[] = [
  {
    element: <Outlet />,
    errorElement: <RootErrorBoundary />,
    children: [
      { index: true, element: <VideoGroups /> },
      { path: "videos", element: <VideoGroups /> },
      { path: "videos/new", element: <NewVideoGroup /> },
      { path: "videos/:id/edit", element: <EditVideoGroup /> },
      { path: "users", element: <Users /> },
      { path: "events", element: <Events /> },
      { path: "events/new", element: <NewEvent /> },
      {
        path: "events/:id",
        element: <EventDetail />,
        children: [{ path: "description/:locale", element: <MarkdownEditorDialog /> }],
      },
      { path: "events/:id/edit", element: <EditEvent /> },
    ],
  },
];
