export default function File({ file }) {
  return (
    <>
      <div>
        <p>{file.id}</p>
        <p>{file.thumbnail}</p>
      </div>
      <div>
        <h1>{file.name}</h1>
        <p>{file.description}</p>
      </div>
      <div>
        <p>{file.size}</p>
        <p>{file.extension}</p>
        <p>{file.upload_date}</p>
        <p>{file.update_date}</p>
      </div>
    </>
  );
}
