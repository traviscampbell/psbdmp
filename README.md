# psbdmp

Unofficial API Wrapper for https://psbdmp.ws written in Go.
Allows searching through pastebin dumps by keyword, domain, email, or date ranges.

### Installation

```bash
$> go get github.com/traviscampbell/psbdmp/...
```

### CLI Examples

Fetch a specific dump by it's ID. (prints it to stdout)
```bash
$> psbdmp -dl f1GH3ySG
```

Search all the dumps for mentions of the domain `github.com`.
By itself like this it only returns a list of dump IDs that matched,
to get the actual dump content you can add the `-fetch` flag.
```bash
$> psbdmp -domain github.com (-fetch)?
```

Download all the pastebin dumps posted in the last 3 days storing them
in `/tmp` by default.
```bash
$> psbdmp -since 3 -fetch
```
