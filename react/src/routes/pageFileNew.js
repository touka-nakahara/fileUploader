import { useState } from "react";
import bcrypt from "bcryptjs";
import styles from "../routes-css/new.module.css";

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

  //RV nakaharaY 途中削除したら際レンダリングされるのでindexがずれないはず...
  function changeMetaDataName(index, name) {
    const updateMetadata = [...metadata];
    const updateElement = metadata[index];

    updateElement.name = name;

    updateMetadata[index] = updateElement;

    setMetadata(updateMetadata);
  }

  //RV nakaharaY もっと効率の良い書き方があるはず
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

    const originalFiles = files;
    const originalMeta = metadata;

    // forループにするか
    files.forEach(async (file, index) => {
      //　パスワードを暗号化する処理
      const salt = bcrypt.genSaltSync(10);

      const password = metadata[index].password;
      let hash = "";
      if (password != "") {
        hash = bcrypt.hashSync(password, salt);
      }

      const formData = new FormData();

      const fileExtension = retriveExtension(file.name);

      const fileData = {
        name: metadata[index].name,
        password: hash,
        description: metadata[index].description,
        extension: fileExtension,
        size: file.size,
      };

      formData.append("file", JSON.stringify(fileData));
      formData.append("data", file);

      const result = await fetch("http://127.0.1:8888/api/files", {
        method: "POST",
        body: formData,
      });

      if (!result.ok) {
        const response = await result.json();
        alert(response.message);
      } else {
        // 成功したらstateから削除
        //RV 非同期に行われるので, indexとの関係を保持できていない
        const removeTarget = originalFiles[index];
        const removeMetadata = originalMeta[index];
        setFiles((prevFiles) => prevFiles.filter((t, _) => t !== removeTarget));
        setMetadata((prevMetadata) =>
          prevMetadata.filter((t, _) => t !== removeMetadata)
        );
      }
    });
  }

  return (
    <div className={styles.fileUplaodNew}>
      <FileUploadField addFile={addFile} handleUpload={handleUpload} />
      <FileSentButton handleUpload={handleUpload} />
      <FilesPreview
        files={files}
        metadata={metadata}
        handleNameChange={changeMetaDataName}
        handlePasswordChange={changeMetaDataPassword}
        handleDescriptionChange={changeMetaDataDescription}
      />
    </div>
  );
}

function FileUploadField({ addFile }) {
  return (
    <>
      <div>
        <p>ファイルをアップロードします</p>
        <FileUploadButton addFile={addFile} />
      </div>
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
      <input
        className={styles.uploadButton}
        type="file"
        accept="*"
        onChange={handleFileChange}
        multiple
      />
    </>
  );
}

//RV nakaharaY これまとめても良い
function FileSentButton({ handleUpload }) {
  return <button onClick={handleUpload}>ファイルを送信</button>;
}

//　ファイルプレビュー配列
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

// ファイルプレビュー
function FilePreview({
  file, // 実はこれいらないんか...
  metadata,
  index,
  handleNameChange,
  handlePasswordChange,
  handleDescriptionChange,
}) {
  return (
    <div className={styles.filePreviewItem}>
      <div>
        <input
          className={styles.inputFieldTitle}
          type="text"
          placeholder="ファイル名"
          name="filename"
          value={metadata.name}
          onChange={(e) => handleNameChange(index, e.target.value)}
        ></input>
      </div>
      <div>
        <p>{file.size} バイト</p>
        <input
          className={styles.inputField}
          type="text"
          placeholder="パスワードを入力してください"
          value={metadata.password}
          onChange={(e) => handlePasswordChange(index, e.target.value)}
        ></input>
        <DescriptionArea
          metadata={metadata}
          index={index}
          handleDescriptionChange={handleDescriptionChange}
        />
      </div>
    </div>
  );
}

// 説明のページs
function DescriptionArea({ metadata, index, handleDescriptionChange }) {
  return (
    <textarea
      className={styles.textArea}
      value={metadata.description}
      onChange={(e) => handleDescriptionChange(index, e.target.value)}
    />
  );
}

// 拡張子を取り除く
function removeExtension(filename) {
  var lastDotIndex = filename.lastIndexOf(".");
  if (lastDotIndex !== -1) {
    return filename.substring(0, lastDotIndex);
  }
  return filename;
}

// 拡張子を取得する
function retriveExtension(filename) {
  var lastDotIndex = filename.lastIndexOf(".");
  if (lastDotIndex === -1) {
    return "";
  }
  var extension = filename.substring(lastDotIndex + 1);
  return extension;
}
