# Xelon SDK for Go

[![Tests](https://github.com/Xelon-AG/xelon-sdk-go/actions/workflows/tests.yaml/badge.svg)](https://github.com/Xelon-AG/xelon-sdk-go/actions)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/Xelon-AG/xelon-sdk-go)

xelon-sdk-go is the official Xelon SDK for the Go programming language.

## Installation

```sh
# X.Y.Z is the version you need
go get github.com/Xelon-AG/xelon-sdk-go@vX.Y.Z


# for non Go modules usage or latest version
go get github.com/Xelon-AG/xelon-sdk-go
```

## Usage

```go
import "github.com/Xelon-AG/xelon-sdk-go"
```

Create a new Xelon client, then use the exposed services to access
different parts of the Xelon API.

### Authentication

Currently, Bearer token is the only method of authenticating with the API.
You can learn how to obtain it [here](https://www.xelon.ch/docs/xelon-api-101#authorize-youself).
Then use your token to create a new client:

```go
client := xelon.NewClient("my-secret-token")
```

If you want to specify more parameters by client initialization, use
`With...` methods and pass via option pattern:

```go
var opts []xelon.ClientOption
opts = append(opts, xelon.WithBaseURL(baseURL))
opts = append(opts, xelon.WithClientID(clientID))
opts = append(opts, xelon.WithUserAgent(userAgent))

client := xelon.NewClient("my-secret-token", opts...)
```

### Examples

List all ssh keys for the user.

```go
func main() {
  client := xelon.NewClient("my-secret-token")

  // list all ssh keys for the authenticated user
  ctx := context.Background()
  sshKeys, _, err := client.SSHKeys.List(ctx)
}
```

Upload a small local ISO directly to a cloud datastore.

```go
file, err := os.Open("oemdrv.iso")
if err != nil {
  // handle error
}
defer file.Close()

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

iso, _, err := client.ISOs.Upload(ctx, &xelon.ISOUploadRequest{
  CategoryID: 1,
  CloudID:    "cloud-123",
  File:       file,
  Filename:   file.Name(),
  Name:       "oemdrv",
})
```

ISO uploads use a two-minute SDK fallback timeout. A caller-provided context
deadline overrides that fallback. When `WithHTTPClient` is used, the supplied
HTTP client owns timeout policy instead. The server documents a 20 MB maximum
for direct uploads; use URL-based `ISOs.Create` for larger installation media.

## Contributing

We love pull requests! Please see the [contribution guidelines](.github/CONTRIBUTING.md).

## License

This SDK is distributed under the Apache-2.0 license, see [LICENSE](LICENSE) for more information.
