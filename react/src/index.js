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
  try {
    const response = await fetch("http://127.0.0.1:8888/api/files");
  } catch (e) {
    console.log("サーバーの応答がありません");
    throw new Error("サーバーの応答がありません");
    //RV nakahara Y throw new Responseが問題だったぽいな 余裕があれば調べる
  }

  if (!response.ok) {
    // サーバーエラーが起きた場合
    throw new Response("サーバーの応答がありません", {
      status: response.status,
    });
  }

  const files = await response.json();

  return { files };
}

//RV nakaharaY fileLoaderがとってきたデータの真部分集合だからこれはデータ的には意味のないFetch
export async function fileLoader({ params }) {
  const { id } = params;
  //TODO マジックナンバーをConfigに
  const response = await fetch(`http://127.0.0.1:8888/api/files/${id}`);

  if (!response.ok) {
    throw new Response("サーバーの応答がありません", {
      status: response.status,
    });
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
        loader: fileLoader,
      },
      {
        path: "/files/new",
        element: <FileNew />,
      },
    ],
  },
]);

const root = createRoot(document.getElementById("root"));
root.render(
  // <StrictMode>
  <RouterProvider router={router} />
  // </StrictMode>
);
