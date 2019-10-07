# crossplane-tools [![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/crossplaneio/crossplane-tools)

Experimental code generators for [Crossplane] controllers.

`angryjet` is the only extant tool within crossplane-tools. It will detect Go
structs that appear to be capable of satisfying crossplane-runtime's interfaces
(such as [`resource.Managed`]) and automatically generate the method set
required to satisfy that interface. A struct is considered capable of satisfying
crossplane-runtime's interfaces based on the heuristics described in the
[Crossplane Services Developer Guide], for example a managed resource must:

* Embed a [`ResourceStatus`] struct in their `Status` struct.
* Embed a [`ResourceSpec`] struct in their `Spec` struct.
* Embed a `Parameters` struct in their `Spec` struct.

Methods are not written if they are already defined outside of the file that
would be generated. Use the `//+crossplane:generate:methods=false` comment
marker to explicitly disable generation of any methods for a type. Use `go
generate` to generate your Crossplane API types by adding a generate marker to
the top level of your `api/` directory, for example:

```go
// Generate crossplane-runtime methodsets (resource.Claim, etc)
//go:generate go run ../vendor/github.com/crossplaneio/crossplane-tools/cmd/angryjet/main.go generate-methodsets ./...
```

```console
$ angryjet generate-methodsets --help
usage: angryjet generate-methodsets [<flags>] [<packages>]

Generate a Crossplane method sets.

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --header-file=HEADER-FILE  The contents of this file will be added to the top of all generated files.
  --filename-managed="zz_generated.managed.go"  
                             The filename of generated managed resource files.
  --filename-claim="zz_generated.claim.go"  
                             The filename of generated resource claim files.
  --filename-portable-class="zz_generated.portableclass.go"  
                             The filename of generated portable class files.
  --filename-portable-class-list="zz_generated.portableclasslist.go"  
                             The filename of generated portable class list files.
  --filename-non-portable-class="zz_generated.nonportableclass.go"  
                             The filename of generated non-portable class files.

Args:
  [<packages>]  Package(s) for which to generate methods, for example github.com/crossplaneio/crossplane/apis/...
```

[Crossplane]: https://crossplane.io
[`resource.Managed`]: https://godoc.org/github.com/crossplaneio/crossplane-runtime/pkg/resource#Managed
[`ResourceSpec`]: https://godoc.org/github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1#ResourceSpec
[`ResourceStatus`]: https://godoc.org/github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1#ResourceStatus
[Crossplane Services Developer Guide]: https://crossplane.io/docs/v0.3/services-developer-guide.html#defining-resource-kinds