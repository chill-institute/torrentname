# Security

`torrentname` is a Go library that parses untrusted release-name strings. It
does not call external services, execute release-name content, or open files
named by parsed input.

## Reporting

Please report suspected vulnerabilities through GitHub Security Advisories for
this repository. If advisories are unavailable, open a minimal public issue that
describes the affected behavior without publishing exploit details.

## Parser Safety

Security-relevant parser changes should keep these properties intact:

- malformed input does not panic
- parsing remains deterministic and dependency-free at runtime
- tests cover minimized crash or hang repro cases
- fuzz coverage is updated when a new input class is added

Run the local safety baseline before proposing a fix:

```bash
mise run verify
mise run test:fuzz
```
