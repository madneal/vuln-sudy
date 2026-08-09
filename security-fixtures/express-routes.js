const express = require("express");
const { exec } = require("node:child_process");
const fs = require("node:fs");
const path = require("node:path");
const sqlite3 = require("sqlite3");

const app = express();
const database = new sqlite3.Database(":memory:");
const publicDirectory = path.resolve(__dirname, "public");

// These are intentionally unsafe routes for CodeQL data-flow testing.
app.get("/command", (req, res) => {
  exec(`ping -c 1 ${req.query.host}`, (error, stdout) => {
    if (error) return res.status(500).send(error.message);
    return res.type("text").send(stdout);
  });
});

app.get("/sql", (req, res) => {
  const query = `SELECT id, username FROM users WHERE username = '${req.query.username}'`;
  database.all(query, (error, rows) => {
    if (error) return res.status(500).send(error.message);
    return res.json(rows);
  });
});

app.get("/file", (req, res) => {
  const filePath = path.join(publicDirectory, req.query.file);
  fs.readFile(filePath, (error, contents) => {
    if (error) return res.status(404).send("Not found");
    return res.type("application/octet-stream").send(contents);
  });
});

app.get("/search", (req, res) => {
  return res
    .type("html")
    .send(`<html><body><h1>Results for ${req.query.q}</h1></body></html>`);
});

// Deliberately no app.listen(): this module is a static-analysis fixture only.
module.exports = app;

