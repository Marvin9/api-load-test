// curl http://localhost:8000 -H "Cache-Control: no-store" -w "@curl-format.txt" -o /dev/null -s
const http = require("http");

/**
 * @type {http.RequestListener}
 */
const server = (req, res) => {
  let j = req.headers['x-http-load-test']
  console.time(`request-${j}`)
  console.log(`request: ${req.url}`);
  res.writeHead(200, {
    "Content-Type": "application/json",
    "Cache-Control": "no-store",
  });
  setTimeout(() => {
    res.end();
    console.timeEnd(`request-${j}`);
  }, 256)
};

http
  .createServer(server)
  .listen(8000)
  .on("listening", () => console.log("server running"));
