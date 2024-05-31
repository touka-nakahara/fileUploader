import styles from "../routes-css/index.module.css";

export default function File({ file }) {
  return (
    <>
      <div className={styles.fileCard}>
        <div>
          <p>{file.thumbnail}</p>
        </div>
        <div>
          <h1>{file.name}</h1>
          <p>{file.description}</p>
        </div>
        <div>
          <p>{Math.floor(file.size / 1000)}KB</p>
          <p>{file.extension}</p>
          <p>{file.update_date}</p>
        </div>
      </div>
    </>
  );
}
