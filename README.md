# angryjet [![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/negz/angryjet)

An experimental code generator for [Crossplane] controllers.

Currently `angryjet` will detect Go structs that appear to be capable of
satisfying crossplane-runtime's [`resource.Managed`] interface and automatically
generate the method set required to satisfy that interface. A struct is
considered capable of satisfying `resource.Managed` if its `Spec` field is a
struct that embeds a [`ResourceSpec`] and its `Status` field is a struct that
embeds a [`ResourceStatus`]. The method set is written to
`zz_generated.managed.go` by default. Methods are not written if they are
already defined.

```
usage: angryjet [<flags>] [<packages>]

Generates Crossplane API type methods.

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --base-dir=/Users/negz/control/go/src  
                             Generated files are written to their package paths relative to this directory.
  --prefix="zz_generated."   This string is prepended to the names of all generated files.
  --header-file=HEADER-FILE  The contents of this file will be added to the top of all generated files.

Args:
  [<packages>]  Package(s) for which to generate Crossplane methods, for example github.com/crossplaneio/crossplane/apis/...
```

[Crossplane]: https://crossplane.io
[`resource.Managed`]: https://godoc.org/github.com/crossplaneio/crossplane-runtime/pkg/resource#Managed
[`ResourceSpec`]: https://godoc.org/github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1#ResourceSpec
[`ResourceStatus`]: https://godoc.org/github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1#ResourceStatus
