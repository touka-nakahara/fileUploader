import { useState } from "react";

export default function FileNew() {
  const [files, setFiles] = useState([]);
  const [metadata, setMetadata] = useState([]);

  function addFile(newFile) {
    setFiles((prevFiles) => [...prevFiles, newFile]);
    setMetadata((prevMetadata) => [
      ...prevMetadata,
      {
        name: removeExtension(newFile.name),
        password: "",
        description: "",
      },
    ]);
  }

  async function handleUpload() {
    if (!files || files.length == 0) {
      return;
    }

    // forループにするか
    const formData = new FormData();
    const fileData = {
      name: files[0].name,
      password: "1234",
      description: "1234",
    };

    formData.append("file", JSON.stringify(fileData));
    formData.append("data", files[0]);

    const result = await fetch("http://127.0.1:8888/api/files", {
      method: "POST",
      body: formData,
    });

    //TODO エラーを起こす方法を知らない...
    if (!result.ok) {
      alert(result);
    }
  }

  return (
    <>
      <FileUploadField addFile={addFile} />
      <FileSentButton handleUpload={handleUpload} />
      <FilesPreview files={files} metadata={metadata} />
    </>
  );
}

function FileUploadField({ addFile }) {
  return (
    <>
      <FileUploadButton addFile={addFile} />
    </>
  );
}

function FileUploadButton({ addFile }) {
  function handleFileChange(e) {
    if (e.target.files) {
      Array.from(e.target.files).forEach((file) => addFile(file));
    }
  }
  return (
    <>
      <div>
        <p>ファイルアップロード</p>
        <input type="file" accept="*" onChange={handleFileChange} multiple />
      </div>
    </>
  );
}

function FileSentButton({ handleUpload }) {
  return (
    <>
      <div>
        <button onClick={handleUpload}>ファイルを送信</button>
      </div>
    </>
  );
}

function FilesPreview({ files, metadata }) {
  if (!files || files.length == 0) {
    return <></>;
  }

  const filePreviews = files.map((file, index) => (
    <li key={index}>
      <FilePreview file={file} metadata={metadata[index]} index={index} />
    </li>
  ));

  return <ul>{filePreviews}</ul>;
}

function FilePreview({ file, metadata, index }) {
  return (
    <div>
      <p>可能ならサムネ</p>
      <input
        type="text"
        placeholder="ファイル名"
        name="filename"
        value={metadata.name}
        onChange={(e) => console.log(e)}
      ></input>
      <p>{file.size} バイト</p>
      <input type="text" placeholder="パスワードを入力してください"></input>
      <input type="text" placeholder="このファイルは..."></input>
    </div>
  );
}

function removeExtension(filename) {
  var lastDotIndex = filename.lastIndexOf(".");
  if (lastDotIndex !== -1) {
    return filename.substring(0, lastDotIndex);
  }
  return filename;
}
