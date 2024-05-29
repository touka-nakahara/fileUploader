import React, { StrictMode } from "react";
import { RouterProvider, Link } from "react-router-dom";
import { createRoot } from "react-dom/client";
import { createBrowserRouter } from "react-router-dom";
import "./styles.css";

import FileIndex from "./routes/pageFileIndex";
import FileDetail from "./routes/pageFileDetail";
import FileNew from "./routes/pageFileNew";
import ErrorPage from "./routes/error";
import Root from "./routes/root";

export async function filesLoader() {
  // ここにFetchを書くってことか

  //TODO マジックナンバーをConfigに
  const response = await fetch("http://127.0.0.1:8888/api/files");

  if (!response.ok) {
    throw new Response("Failed to fetch data", { status: response.status });
  }

  const files = await response.json();

  return { files };
}

// fileLoaderが撮ってきたデータの真部分集合だからこれは本質的には意味のないFetch
export async function fileLoader({ params }) {
  // ここにFetchを書くってことか
  const { id } = params;
  //TODO マジックナンバーをConfigに
  const response = await fetch(`http://127.0.0.1:8888/api/files/${id}`);

  if (!response.ok) {
    throw new Response("Failed to fetch data", { status: response.status });
  }

  const file = await response.json();

  return { file };
}

// ルーティング
const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "/",
        element: <FileIndex />,
        loader: filesLoader,
      },
      {
        path: "/files/:id",
        element: <FileDetail />,
        errorElement: <ErrorPage />,
        loader: fileLoader,
      },
      {
        path: "/files/new",
        element: <FileNew />,
        errorElement: <ErrorPage />,
      },
    ],
  },
]);

const root = createRoot(document.getElementById("root"));
root.render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>
);
