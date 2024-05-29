import { Outlet, Link } from "react-router-dom";

import styles from "../routes-css/root.module.css";

export default function Header() {
  return (
    <>
      <div className={styles.header}>
        <RootSender />
        <NewFileSender />
      </div>
      <Outlet />
    </>
  );
}

function RootSender() {
  return (
    <>
      <Link to={"/"} className={styles.title}>
        FileUploader
      </Link>
    </>
  );
}

function NewFileSender() {
  return (
    <>
      <Link to={"files/new"} className={styles.upload}>
        Upload
      </Link>
    </>
  );
}
