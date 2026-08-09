# GitHub vulnerability-scanning fixtures

This directory contains intentionally insecure JavaScript examples for testing
GitHub CodeQL and Dependabot detection. They are not application code and must
not be deployed, exposed to a network, or copied into a production service.

Expected findings include command injection, SQL injection, path traversal,
reflected cross-site scripting, and vulnerable npm dependencies. The
`express-routes.js` and `go/routes.go` files register realistic HTTP handlers,
but deliberately do not start a server. The files are kept separate from the
repository root so that their purpose is unambiguous.

The CodeQL workflow uses the `security-extended` query suite and analyzes both
JavaScript and Go. The Go fixture is a standalone module so CodeQL can build it
without fetching application dependencies. The `/go/complex/*` routes add
multi-hop flows through decoders, value objects, interfaces, service methods,
and output builders before reaching the vulnerable sinks.
