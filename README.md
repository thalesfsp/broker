# broker

`broker` allows to register `Listener`s to listen for events. A `broker` broadcast an event to all `Listener`s, 1:N. Broker's goal isn't to broadcast events cross-services, for this case use NATS, Redis, and etc. Goal is to propagate an event - a change, internally and react to it.

## Install

`$ go get github.com/thalesfsp/broker`

### Specific version

Example: `$ go get github.com/thalesfsp/broker@v1.2.3`

## Usage

See [`example_test.go`](example_test.go) file.

### Documentation

Run `$ make doc` or check out [online](https://pkg.go.dev/github.com/thalesfsp/broker).

## Development

Check out [CONTRIBUTION](CONTRIBUTION.md).

### Release

1. Update [CHANGELOG](CHANGELOG.md) accordingly.
2. Once changes from MR are merged.
3. Tag and release.

## Roadmap

Check out [CHANGELOG](CHANGELOG.md).
