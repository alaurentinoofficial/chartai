import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Layout from "@/components/layout";
import { ChartsPage } from "@/pages/charts";
import { DatabasesPage } from "@/pages/databases";

const router = createBrowserRouter([
  {
    path: "/",
    children: [
      {
        path: "/",
        element: <Layout>Home</Layout>,
      },
      {
        path: "charts",
        element: <Layout><ChartsPage /></Layout>,
      },
      {
        path: "databases",
        element: <Layout><DatabasesPage /></Layout>,
      },
      {
        path: "settings",
        element: <Layout>Settings</Layout>,
      },
	  {
		  path: "*",
		  element: <Layout>Not Found</Layout>,
	  }
    ],
  },
]);

export function Routers() {
  return <RouterProvider router={router} />;
} 
