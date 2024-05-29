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

  //TODO 途中削除したら際レンダリングされるのでindexがずれないはず...
  function changeMetaDataName(index, name) {
    const updateMetadata = [...metadata];
    const updateElement = metadata[index];

    updateElement.name = name;

    updateMetadata[index] = updateElement;

    setMetadata(updateMetadata);
  }

  // アホすぎて笑えるわ
  function changeMetaDataPassword(index, password) {
    const updateMetadata = [...metadata];
    const updateElement = metadata[index];

    updateElement.password = password;

    updateMetadata[index] = updateElement;

    setMetadata(updateMetadata);
  }

  function changeMetaDataDescription(index, description) {
    const updateMetadata = [...metadata];
    const updateElement = metadata[index];

    updateElement.description = description;

    updateMetadata[index] = updateElement;

    setMetadata(updateMetadata);
  }

  async function handleUpload() {
    if (!files || files.length == 0) {
      return;
    }

    // forループにするか
    files.forEach(async (file, index) => {
      //　パスワードを暗号化する処理

      const formData = new FormData();

      const fileData = {
        name: metadata.name,
        password: metadata[index].password,
        description: metadata[index].description,
        extension: retriveExtension(file.name),
      };

      formData.append("file", JSON.stringify(fileData));
      formData.append("data", file);

      const result = await fetch("http://127.0.1:8888/api/files", {
        method: "POST",
        body: formData,
      });

      //TODO エラーを起こす方法を知らない...
      if (!result.ok) {
        alert(result);
      }
    });
  }

  return (
    <>
      <FileUploadField addFile={addFile} />
      <FileSentButton handleUpload={handleUpload} />
      <FilesPreview
        files={files}
        metadata={metadata}
        handleNameChange={changeMetaDataName}
        handlePasswordChange={changeMetaDataPassword}
        handleDescriptionChange={changeMetaDataDescription}
      />
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

// マジで終わってるどうしたらええんや？
function FilesPreview({
  files,
  metadata,
  handleNameChange,
  handlePasswordChange,
  handleDescriptionChange,
}) {
  if (!files || files.length == 0) {
    return <></>;
  }

  const filePreviews = files.map((file, index) => (
    <li key={index}>
      <FilePreview
        file={file}
        metadata={metadata[index]}
        index={index}
        handleNameChange={handleNameChange}
        handlePasswordChange={handlePasswordChange}
        handleDescriptionChange={handleDescriptionChange}
      />
    </li>
  ));

  return <ul>{filePreviews}</ul>;
}

function FilePreview({
  file, // 実はこれいらないんか...
  metadata,
  index,
  handleNameChange,
  handlePasswordChange,
  handleDescriptionChange,
}) {
  return (
    <div>
      <p>可能ならサムネ</p>
      <input
        type="text"
        placeholder="ファイル名"
        name="filename"
        value={metadata.name}
        onChange={(e) => handleNameChange(index, e.target.value)}
      ></input>
      <p>{file.size} バイト</p>
      <input
        type="text"
        placeholder="パスワードを入力してください"
        value={metadata.password}
        onChange={(e) => handlePasswordChange(index, e.target.value)}
      ></input>
      <input type="text" placeholder="このファイルは..."></input>
      <DescriptionArea
        metadata={metadata}
        index={index}
        handleDescriptionChange={handleDescriptionChange}
      />
    </div>
  );
}

function DescriptionArea({ metadata, index, handleDescriptionChange }) {
  return (
    <textarea
      value={metadata.description}
      onChange={(e) => handleDescriptionChange(index, e.target.value)}
    />
  );
}

function removeExtension(filename) {
  var lastDotIndex = filename.lastIndexOf(".");
  if (lastDotIndex !== -1) {
    return filename.substring(0, lastDotIndex);
  }
  return filename;
}

function retriveExtension(filename) {
  var lastDotIndex = filename.lastIndexOf(".");
  if (lastDotIndex === -1) {
    return "";
  }
  var extension = filename.substring(lastDotIndex + 1);
  return extension;
}
