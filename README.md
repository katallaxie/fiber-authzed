# 👮 Authzed

## OpenAPI

Using [OpenAPI Extensions](https://swagger.io/docs/specification/openapi-extensions/) individual operations can be with [SpiceDB](https://authzed.com/spicedb).

```yaml
x-fiber-authzed:
  permission: view
  subject:
    type: user
    components:
      - in: params
        name: teamId
  object:
    type: document
    components:
      - in: params
        name: documentId
```

## Examples

See [examples](https://github.com/katallaxie/fiber-authzed/tree/main/examples) to understand the provided interfaces.

## Development

```bash
zed context set dev localhost:50051 example --insecure
```

```bash
zed import ./examples/authzed-download-2233c2.yaml
```

## License

[MIT](/LICENSE)