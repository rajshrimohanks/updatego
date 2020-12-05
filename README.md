# updatego

An absolutely unnecessary utility to install/update Go in **my** machine. It probably won't work for you unless you use a similar setup as mine.

- [updatego](#updatego)
  - [Pre-requisites](#pre-requisites)
  - [Build](#build)
  - [Install](#install)
  - [Usage](#usage)
  - [Perks](#perks)
  - [Why is this needed?](#why-is-this-needed)

## Pre-requisites

For this to work, you need to have Go install to in `/usr/local`. You also should have the following lines added in your `~/.bashrc`:

```bash
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$(go env GOPATH)/bin
```

This will let you take advantage of using Go modules to manage packages and your Go projects.

## Build

Just clone and run `go build`.

## Install

If you don't want to build and manage your binary yourself, you can just use the below commands to directly setup the tool in your system:

```bash
go get -u github.com/rajshrimohanks/updatego
sudo mv $(go env GOPATH)/bin/updatego /usr/local/bin/
```

Now you can just do `sudo updatego` to run the utility.

## Usage

You need to run the utility with `sudo` since it is setting up your Go inside `/usr/local`.

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

## Perks

Since all we are doing is just deleting the existing go setup and extracting the new one, we can use this to roll back go versions as well (_as if this wasn't obvious, duh!_)

```bash
sudo updatego -version 1.14.0
```

## Why is this needed?

I totally have no idea. You actually need Go to build this in the first place. You could also use `go get` to pull in a new version of Go. But `go get` will only bring your new version as a package. It won't replace your primary install. This utility will do that. So I guess, that's a usecase? :P
