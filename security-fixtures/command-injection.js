const { exec } = require("node:child_process");

// INTENTIONAL VULNERABILITY: user input reaches a shell command unchanged.
function pingHost(req, res) {
  const host = req.query.host;

  exec(`ping -c 1 ${host}`, (error, stdout) => {
    if (error) {
      return res.status(500).send(error.message);
    }

    return res.type("text").send(stdout);
  });
}

module.exports = { pingHost };

