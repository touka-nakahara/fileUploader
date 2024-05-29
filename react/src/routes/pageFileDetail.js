import File from "./file";
import { useLoaderData } from "react-router-dom";

export default function FileDetail() {
  const { file } = useLoaderData();
  return (
    <div>
      <File file={file.data} />
    </div>
  );
}
