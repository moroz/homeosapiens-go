import { Outlet, type RouteObject } from "react-router";

import { RootErrorBoundary } from "./root";
import Events from "./routes/events/events";
import { Users } from "./routes/users";
import VideoGroups from "./routes/video-groups/video-groups";
import NewVideoGroup from "./routes/video-groups/new-video-group";
import { EditVideoGroup } from "./routes/video-groups/edit-video-group";
import { MarkdownEditorDialog } from "~/components/markdown-editor";
import { UserDetails } from "./routes/users/user-details";
import { NewBlogPost, BlogPostDetails, BlogPosts, EditBlogPost } from "~/routes/blog-posts";
import { Products } from "~/routes/products";
import { Orders, OrderDetails } from "~/routes/orders";
import {
  EditEvent,
  EventAttendants,
  NewEvent,
  EventDetails,
  EnrollStudentDialog,
} from "~/routes/events";

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
      { path: "users/:id", element: <UserDetails /> },
      { path: "events", element: <Events /> },
      { path: "products", element: <Products /> },
      { path: "orders", element: <Orders /> },
      { path: "orders/:id", element: <OrderDetails /> },
      { path: "blog-posts", element: <BlogPosts /> },
      { path: "blog-posts/new", element: <NewBlogPost /> },
      { path: "blog-posts/:id", element: <BlogPostDetails /> },
      { path: "blog-posts/:id/edit", element: <EditBlogPost /> },
      { path: "events/new", element: <NewEvent /> },
      {
        path: "events/:id",
        element: <EventDetails />,
        children: [{ path: "description/:locale", element: <MarkdownEditorDialog /> }],
      },
      { path: "events/:id/edit", element: <EditEvent /> },
      {
        path: "events/:id/attendants",
        element: <EventAttendants />,
        children: [{ path: "add", element: <EnrollStudentDialog /> }],
      },
    ],
  },
];
