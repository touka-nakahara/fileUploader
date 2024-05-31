import { useLoaderData, useNavigate, Link } from "react-router-dom";
import File from "./file.js";
import { useEffect, useState } from "react";

export default function FileIndex() {
  const { files } = useLoaderData();
  const [filesData, setFilesData] = useState(null);
  const [pageIndex, setPageIndex] = useState(1);

  async function handleDelete({ file, password }) {
    const result = await fetch(`http://127.0.1:8888/api/files/${file.id}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ password: password }),
    });

    if (!result.ok) {
      const message = await result.json();
      alert(message.message);
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
        <Page
          pageIndex={pageIndex}
          setPageIndex={setPageIndex}
          setFilesData={setFilesData}
        />
        <Search setFilesData={setFilesData} />
        <Sort setFilesData={setFilesData} pageIndex={pageIndex} />
        <FileTable files={filesData} handleDelete={handleDelete} />
      </div>
    </>
  );
}

function FileTable({ files, handleDelete }) {
  if (!files || files.len == 0) {
    return <></>;
  }
  const fileItems = files.map((file) => {
    return (
      <li key={file.id}>
        <FileRow file={file} />
        <FileUtiles file={file} handleDelete={handleDelete} />
      </li>
    );
  });
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

// 消したらリロードさせればよかったか？
export function FileUtiles({ file, handleDelete }) {
  const [password, setPassword] = useState("");

  let passwordInput = <></>;
  if (file.has_password) {
    passwordInput = (
      <input
        type="text"
        value={password}
        onChange={(e) => {
          setPassword(e.target.value);
        }}
        placeholder="input file's password"
      />
    );
  }

  async function handleDownload() {
    // パスワードがあるならパスワードをハッシュにしなくていいんかな...わかんない

    const result = await fetch(
      `http://127.0.1:8888/api/files/${file.id}/download`,
      {
        method: "post",
        mode: "cors",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ password: password }),
      }
    );

    if (result.status != 200) {
      const message = await result.json();
      alert(message.message);
      return;
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
    tempLink.download = file.name + "." + file.extension;
    tempLink.click();
  }

  return (
    <>
      {passwordInput}
      <button onClick={handleDownload}>dl</button>
      <button onClick={() => handleDelete({ file, password })}>delete</button>
    </>
  );
}

function Search({ setFilesData }) {
  const [searchWord, setSerachWord] = useState("");

  // 検索を飛ばす
  async function handelSearch({ searchWord }) {
    const response = await fetch(
      `http://127.0.0.1:8888/api/files?search=${searchWord}`
    );

    if (!response.ok) {
      throw new Response("Failed to fetch data", { status: response.status });
    }

    const file = await response.json();

    setFilesData(file.data);
  }

  return (
    <>
      <input
        type="search"
        value={searchWord}
        placeholder="Search..."
        onChange={(e) => {
          setSerachWord(e.target.value);
        }}
      />
      <button
        onClick={() => {
          handelSearch({ searchWord });
        }}
      >
        Search
      </button>
    </>
  );
}

function Sort({ setFilesData, pageIndex }) {
  // DESC, ASC
  const [nameOrder, setNameOrder] = useState(true);
  const [updateDataOrder, setUpdateDataOrder] = useState(true);
  const [sizeOrder, setSizeOrder] = useState(true);

  // 検索を飛ばす
  async function handelSort(sortTag, order) {
    let query = `http://127.0.0.1:8888/api/files?sort=${sortTag}&page=${pageIndex}`;

    // TrueだったらDESC
    if (order) {
      query += `&ordered=desc`;
    } else {
      query += `&ordered=asc`;
    }

    const response = await fetch(query);

    if (!response.ok) {
      throw new Response("Failed to fetch data", { status: response.status });
    }

    const file = await response.json();

    setFilesData(file.data);
  }

  return (
    <div>
      <button
        onClick={() => {
          //RV こけたらどうするん？
          handelSort("name", nameOrder);
          setNameOrder(!nameOrder);
        }}
      >
        name
      </button>
      <button
        onClick={() => {
          handelSort("update_date", updateDataOrder);
          setUpdateDataOrder(!updateDataOrder);
        }}
      >
        last_updated
      </button>
      <button
        onClick={() => {
          handelSort("size", sizeOrder);
          setSizeOrder(!sizeOrder);
        }}
      >
        size
      </button>
    </div>
  );
}

function Page({ pageIndex, setPageIndex, setFilesData }) {
  // 次がなくてもあってもボタンを置いておく
  //RV クエリをクリックするたびに飛ばすのが正しい形なのだろうか？
  async function handlePage(index) {
    if (index < 1) {
      // なにもしない
      return;
    }

    const response = await fetch(
      `http://127.0.0.1:8888/api/files?page=${index}`
    );

    if (!response.ok) {
      throw new Response("Failed to fetch data", { status: response.status });
    }

    const file = await response.json();

    // これはどうなんだ？
    if (file.data) {
      setFilesData(file.data);
      setPageIndex(index);
    }
  }
  return (
    <>
      <button
        onClick={() => {
          handlePage(pageIndex - 1);
        }}
      >
        前
      </button>
      <span>{pageIndex}</span>
      <button
        onClick={() => {
          handlePage(pageIndex + 1);
        }}
      >
        次
      </button>
    </>
  );
}
