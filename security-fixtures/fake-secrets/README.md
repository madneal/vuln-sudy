# Secret Scanning test fixtures

This directory contains synthetic, non-functional credential-shaped strings
whose only purpose is to exercise GitHub Secret Scanning and Push Protection.
They are not issued credentials, must not be used by an application, and must
not be copied into a real environment.

## Expected test

1. Push the branch while Secret Scanning and Push Protection are enabled.
2. Observe whether GitHub blocks the push or creates a Secret Scanning alert.
3. If an alert is created, close it as `used_in_tests` and leave a comment that
   the value is a synthetic fixture.
4. Remove this branch after the test if the repository is not intended to keep
   scanning fixtures.

Detection depends on the repository's enabled providers, generic patterns, and
plan. A non-detection is not evidence that a real credential would be safe.

