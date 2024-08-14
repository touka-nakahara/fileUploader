import http from "k6/http";
import { sleep } from "k6";
import { randomString } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";

export const options = {
  // A number specifying the number of VUs to run concurrently.
  vus: 20,
  // A string specifying the total duration of the test run.
  duration: "60s",
};

const fileList = [
  "fakefile_50.dummy",
  "fakefile_100.dummy",
  "fakefile_450.dummy",
  // "fakefile_500.dummy",
  "image_4k.jpg",
  "image.png",
];

const fileListBlob = [
  open("./sample/" + "fakefile_50.dummy", "b"),
  open("./sample/" + "fakefile_100.dummy", "b"),
  open("./sample/" + "fakefile_450.dummy", "b"),
  // open("./sample/" + "fakefile_500.dummy", "b"),
  open("./sample/" + "image_4k.jpg", "b"),
  open("./sample/" + "image.png", "b"),
];

const downloadList = [];

function uploadDownload() {
  http.get("http://localhost:8888/api/files"); // サイトアクセス

  // ランダムな数字を取得
  const fileNum = Math.floor(Math.random() * (fileList.length - 1));
  const binFile = fileListBlob[fileNum];
  // 適当な文字列8~12文字
  const passwordLength = Math.floor(Math.random() * (12 - 8) + 8);
  const uploadPassword = randomString(passwordLength);

  // 適当な文字列10~100文字
  const descriptionLength = Math.floor(Math.random() * (100 - 10) + 10);
  const description = randomString(descriptionLength);

  const fileData = {
    name: fileList[fileNum],
    password: uploadPassword,
    description: description,
    extension: retriveExtension(fileList[fileNum]),
    size: binFile.byteLength,
  };
  const uploadBody = {
    file: JSON.stringify(fileData),
    data: http.file(binFile),
  };

  // アップロード
  console.log("upload", fileList[fileNum]);
  http.post("http://localhost:8888/api/files", uploadBody); // アップロード
  // const uploadRes = http.post("http://localhost:8888/api/files", uploadBody); // アップロード
  // console.log(uploadRes.json());

  // // 成功していた場合
  // if (uploadRes.status == 200) {
  //   downloadList.push({
  //     id: uploadRes.json().data.id,
  //     password: uploadPassword,
  //   });
  // } else {
  //   return;
  // }

  // console.log(downloadList);

  // // ダウンロード
  // // ランダムなものを選ぶ
  // const dlNum = Math.floor(Math.random() * (downloadList.length - 1));
  // // パスワードを詰める
  // const donwloadPassword = downloadList[dlNum].password;

  // const downloadBody = JSON.stringify({ password: donwloadPassword });

  // const donwloadID = downloadList[dlNum].id;

  // // クエリ発行
  // http.post(
  //   `http://localhost:8888/api/files/${donwloadID}/download`,
  //   downloadBody
  // );

  // console.log("download end");
}

export default function () {
  uploadDownload();
}

function retriveExtension(filename) {
  var lastDotIndex = filename.lastIndexOf(".");
  if (lastDotIndex === -1) {
    return "";
  }
  var extension = filename.substring(lastDotIndex + 1);
  return extension;
}
