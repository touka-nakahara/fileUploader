import { useLoaderData, useNavigate, Link } from "react-router-dom";
import File from "./file.js";
import { useEffect, useState } from "react";

export default function FileIndex() {
  const { files } = useLoaderData();
  const [filesData, setFilesData] = useState(null);

  async function handleDelete({ file }) {
    const password = 1234;

    const result = await fetch(`http://127.0.1:8888/api/files/${file.id}`, {
      method: "POST",
      body: { password: password },
    });

    if (!result.ok) {
      console.log(result);
    } else {
      setFilesData((prevFilesData) => {
        const newFileData = prevFilesData.filter((f) => f.id !== file.id);
        return newFileData;
      });
    }
  }

  // なにこれ...
  // https://ja.react.dev/learn/synchronizing-with-effects
  useEffect(() => {
    setFilesData(files.data);
  }, [files.data]);

  return (
    <>
      <div>
        <SearchBar />
        <div>
          <TableBar />
          <FilterBar />
        </div>
        <FileTable files={filesData} handleDelete={handleDelete} />
        <PageSentBar />
      </div>
    </>
  );
}

function TableBar() {
  //TODO こいつクリックしたらソートか？
  return (
    <div>
      <p>ファイル名</p>
      <p>ファイルサイズ</p>
      <p>拡張子</p>
      <p>最終アップデート日</p>
      <p>アップロード日</p>
      <p>パスワード</p>
    </div>
  );
}

function FilterBar() {
  return (
    <div>
      <p>ソート</p>
      <p>フィルター</p>
    </div>
  );
}

function SearchBar() {
  return (
    <div>
      <p>検索</p>
    </div>
  );
}

//TODO ページ送り
function PageSentBar() {}

function FileTable({ files, handleDelete }) {
  if (!files || files.len == 0) {
    return <></>;
  }
  const fileItems = files.map((file) => (
    <li key={file.id}>
      <FileRow file={file} />
      <DownloadButton file={file} />
      <DeleteButton file={file} handleDelete={handleDelete} />
    </li>
  ));
  return <ul>{fileItems}</ul>;
}

function FileRow({ file }) {
  //TODO BLOBをwebpに変換する
  const navigate = useNavigate();

  const handleClick = () => {
    navigate(`/files/${file.id}`);
  };

  return (
    <>
      <div onClick={handleClick} style={{ cursor: "pointer" }}>
        <File file={file} />
      </div>
    </>
  );
}

function DownloadButton({ file }) {
  const password = 1234;

  async function handleDownload() {
    const result = await fetch(
      `http://127.0.1:8888/api/files/${file.id}/download`,
      {
        method: "POST",
        body: { password: password },
      }
    );

    if (!result.ok) {
      console.log(result);
    }

    // Blobが帰ってくるので
    const responseJson = await result.json();

    //TODO これよくわかってないけど動いてる
    const fileData = Uint8Array.from(atob(responseJson.data.data), (c) =>
      c.charCodeAt(0)
    ); // base64からバイト配列に変換
    const blob = new Blob([fileData]);

    const downloadURL = window.URL.createObjectURL(blob);
    const tempLink = document.createElement("a");
    tempLink.href = downloadURL;
    tempLink.download = file.name + file.extension;
    tempLink.click();
  }

  return <button onClick={handleDownload}>dl</button>;
}

function DeleteButton({ file, handleDelete }) {
  return <button onClick={() => handleDelete({ file })}>delete</button>;
}
