// INTENTIONAL VULNERABILITY: request data is interpolated into an HTML response
// without output encoding.
function search(req, res) {
  const term = req.query.q || "";

  return res
    .type("html")
    .send(`<html><body><h1>Results for ${term}</h1></body></html>`);
}

module.exports = { search };

