# DGSIS Website

Engineering software that lasts.

## Development

Run:

```bash
go run ./cmd/web
```

`make check` runs gofmt, `go vet`, and `go test ./...`.

## CI/CD

GitHub Actions on a self-hosted runner (`runs-on: self-hosted`) registered to this repository. CI runs on pull requests and on pushes to `main` and `develop`. Deploy builds the image on the VPS and switches Traefik blue/green on tag `v*` or `workflow_dispatch`. No registry push and no GitHub Actions secrets.

See [`deploy/README.md`](deploy/README.md).
