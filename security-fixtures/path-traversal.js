const fs = require("node:fs");
const path = require("node:path");

const publicDirectory = path.resolve(__dirname, "public");

// INTENTIONAL VULNERABILITY: a relative path from the request is not confined
// to publicDirectory before the file is opened.
function downloadFile(req, res) {
  const requestedFile = req.query.file;
  const filePath = path.join(publicDirectory, requestedFile);

  fs.readFile(filePath, (error, contents) => {
    if (error) {
      return res.status(404).send("Not found");
    }

    return res.type("application/octet-stream").send(contents);
  });
}

module.exports = { downloadFile };

