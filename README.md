[![Go Documentation](http://godoc.org/github.com/raff/goble?status.svg)](http://godoc.org/github.com/raff/goble)
[![Go Report Card](https://goreportcard.com/badge/github.com/raff/goble)](https://goreportcard.com/report/github.com/raff/goble)
<!--[![Actions Status](https://github.com/raff/goble/workflows/Go/badge.svg)](https://github.com/raff/goble/actions)-->

goble
=====

Go implementation of Bluetooth LE support for OSX (derived from noble/bleno)

This is a port of nodejs [noble](https://github.com/sandeepmistry/noble)/[bleno](https://github.com/sandeepmistry/bleno) for OSX only.

Once I have something working it can maybe integrated with [github.com/paypal/gatt](https://github.com/paypal/gatt), that right now is Linux only.

## Installation

    $ go get github.com/raff/goble
    
## Documentation
http://godoc.org/github.com/raff/goble

## Examples
* examples/main.go : an example of how to use most of the APIs
* examples/discoverer.go : a port of nodejs noble "advertisement-discovery.js" example
* examples/explorer.go : a port of nodejs noble "peripheral-explorer.js" example
