# Xelon SDK for Go

[![Tests](https://github.com/Xelon-AG/xelon-sdk-go/actions/workflows/tests.yaml/badge.svg)](https://github.com/Xelon-AG/xelon-sdk-go/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Xelon-AG/xelon-sdk-go)](https://goreportcard.com/report/github.com/Xelon-AG/xelon-sdk-go)
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

## Contributing

We love pull requests! Please see the [contribution guidelines](.github/CONTRIBUTING.md).

## License

This SDK is distributed under the Apache-2.0 license, see [LICENSE](LICENSE) for more information.
