# crossplane-tools [![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/crossplaneio/crossplane-tools) [![Go Report Card](https://goreportcard.com/badge/github.com/crossplaneio/crossplane-tools)](https://goreportcard.com/report/github.com/crossplaneio/crossplane-tools)

Crossplane Tools for [Crossplane] controllers. See the [Crossplane Services Developer Guide] for more information.

| Tool | Description
| ---  | ---
| [AngryJet] | AngryJet is a code generator for Crossplane controllers. It will detect Go structs that appear to be capable of satisfying crossplane-runtime's interfaces (such as [`resource.Managed`]) and automatically generate the method set required to satisfy that interface.

[AngryJet]: cmd/angryjet/README
[Crossplane]: https://crossplane.io
[Crossplane Services Developer Guide]: https://crossplane.io/docs/v0.3/services-developer-guide.html#defining-resource-kinds
[`resource.Managed`]: https://godoc.org/github.com/crossplaneio/crossplane-runtime/pkg/resource#Managed
