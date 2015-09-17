var url = "http://localhost:3000/sse";

var es = new EventSource(url);
es.addEventListener("open", function (e) {
  console.log("sse:open:", e);
});
es.addEventListener("message", function (e) {
  console.log("sse:message:", e.data);
});
es.addEventListener("error", function (e) {
  console.error("sse:error:", e);
});
