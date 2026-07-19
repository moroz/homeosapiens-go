import type { RouteObject } from "react-router";

import { RootLayout, RootErrorBoundary } from "./root";
import Home from "./routes/home";

export const routes: RouteObject[] = [
  {
    element: <RootLayout />,
    errorElement: <RootErrorBoundary />,
    children: [{ index: true, element: <Home /> }],
  },
];
