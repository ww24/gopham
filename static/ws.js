var url = "ws://localhost:3000/subscribe";

var ws = new WebSocket(url);
ws.addEventListener("open", function (e) {
  console.log("ws:open:", e);
});
ws.addEventListener("close", function (e) {
  console.log("ws:close:", e);
});
ws.addEventListener("message", function (e) {
  console.log("ws:message:", e.data);
});
ws.addEventListener("error", function (e) {
  console.error("ws:error", e);
});

function send(data) {
  var json_str = JSON.stringify({
    channel: "default",
    ttl: 10,
    data: data
  });
  ws.send(json_str);
}
