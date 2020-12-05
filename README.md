# updatego

An absolutely unnecessary utility to install/update Go in **my** machine. It probably won't work for you unless you use a similar setup as mine.

## Pre-requisites

For this to work, you need to have Go install to in `/usr/local`. You also should have the following lines added in your `~/.bashrc`:

```bash
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$(go env GOPATH)/bin
```

This will let you take advantage of using Go modules to manage packages and your Go projects.

## Build

Just clone and run `go build`.

## Usage

You need to run the utility with `sudo` since it is setting up your go inside `/usr/local`.

```bash
sudo updatego -version 1.15.5
```

All options:

```plaintext
  -go-dir string
        the directory inside which the go archive should be extracted (default "/usr/local")
  -version string
        the go version to fetch (default "1.15.5")
```
