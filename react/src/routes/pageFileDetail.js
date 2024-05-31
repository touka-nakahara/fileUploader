import File from "./file";
import { useLoaderData } from "react-router-dom";
import { FileUtiles } from "./pageFileIndex";

export default function FileDetail() {
  const { file } = useLoaderData();

  async function handleDelete({ file, password }) {
    const result = await fetch(`http://127.0.1:8888/api/files/${file.id}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ password: password }),
    });
  }

  return (
    <div>
      <File file={file.data} />
      <FileUtiles file={file.data} handleDelete={handleDelete} />
    </div>
  );
}
