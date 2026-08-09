# GitHub vulnerability-scanning fixtures

This directory contains intentionally insecure JavaScript examples for testing
GitHub CodeQL and Dependabot detection. They are not application code and must
not be deployed, exposed to a network, or copied into a production service.

Expected findings include command injection, SQL injection, path traversal,
reflected cross-site scripting, and vulnerable npm dependencies. The files are
kept separate from the repository root so that their purpose is unambiguous.

