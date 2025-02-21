const { createProxyMiddleware } = require("http-proxy-middleware");

// change this to the address exposing the participant-api
const target = "http://localhost:3231";

module.exports = function (app) {
  app.use(
    "/api",
    createProxyMiddleware({
      target: target,
      secure: false,
      pathRewrite: {
        "^/api/": "/"
      }
    })
  );
};