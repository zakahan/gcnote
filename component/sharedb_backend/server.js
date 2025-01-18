const http = require("http");
const WebSocket = require("ws");
const ShareDB = require("sharedb");
const WebSocketJSONStream = require("@teamwork/websocket-json-stream");
const ShareDBMongo = require("sharedb-mongo");

// 连接 MongoDB 作为 ShareDB 的数据库
const mongoUrl = "mongodb://localhost:27017/gcnote"; // 修改为你的 MongoDB 地址
const db = ShareDBMongo(mongoUrl);
const shareDB = new ShareDB({ db });

// 创建 HTTP 服务器
const server = http.createServer();
const wss = new WebSocket.Server({ server });

wss.on("connection", (ws) => {
    const stream = new WebSocketJSONStream(ws);
    shareDB.listen(stream);
});

// 监听端口 8096
server.listen(8096, () => {
    console.log("ShareDB 服务器运行在 http://localhost:8096，MongoDB 已启用持久化存储");
});
