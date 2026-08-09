const sqlite3 = require("sqlite3");

const database = new sqlite3.Database(":memory:");

// INTENTIONAL VULNERABILITY: the query is built with untrusted input.
function findUser(req, res) {
  const username = req.query.username;
  const query = `SELECT id, username FROM users WHERE username = '${username}'`;

  database.all(query, (error, rows) => {
    if (error) {
      return res.status(500).send(error.message);
    }

    return res.json(rows);
  });
}

module.exports = { findUser };

