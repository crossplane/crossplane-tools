# crossplane-tools [![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/crossplane/crossplane-tools)

Code generators for [Crossplane] controllers.

## angryjet

`angryjet` will detect Go structs that appear to be capable of satisfying
crossplane-runtime's interfaces (such as [`resource.Managed`]) and automatically
generate the method set required to satisfy that interface. A struct is
considered capable of satisfying crossplane-runtime's interfaces based on the
heuristics described in the [Provider Development Guide], for example a
managed resource must:

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
//go:generate go run ../vendor/github.com/crossplane/crossplane-tools/cmd/angryjet/main.go generate-methodsets ./...
```

### Reference Resolvers

In addition to functions that satisfy `resource.Managed`, you can use `angryjet`
to generate a `ResolveReferences` method as well. In order to generate a resolution
call for given field, you need to add the following comment marker:
```
// +crossplane:generate:reference:type=<target type>
```

`<target type>` could either be just the type name of the CRD if it is in the same
package or `<package path>.<target type>` if it is in a different package, such
as `github.com/crossplane/provider-aws/apis/ec2/v1beta1.VPC`.

The generated resolver will use the external name annotation of the target resource
to fetch the value and it assumes that reference field is named as
`FieldNameRef`/`FieldNameRefs if array` and selector field is named as 
`FieldNameSelector`. You can override these defaults by adding the optional comment
markers, see the following example:
```go
type SomeParameters struct {
    // +crossplane:generate:reference:type=github.com/crossplane/provider-aws/apis/ec2/v1beta1.Subnet
    // +crossplane:generate:reference:extractor=github.com/crossplane/provider-aws/apis/ec2/v1beta1.SubnetARN()
    // +crossplane:generate:reference:refFieldName=SubnetIDRefs
    // +crossplane:generate:reference:selectorFieldName=SubnetIDSelector
    SubnetIDs []string `json:"subnetIds,omitempty"`
    
    SubnetIDRefs []xpv1.Reference `json:"subnetIdRefs,omitempty"`
    
    SubnetIDSelector *xpv1.Selector `json:"subnetIdSelector,omitempty"`
}
```

Note that it doesn't make any change to the CRD struct; authors still need to
add `FieldNameRef` and `FieldNameSelector` fields on their own for the generated
code to compile.

### Usage

```console
$ angryjet generate-methodsets --help
usage: angryjet generate-methodsets [<flags>] [<packages>]

Generate a Crossplane method sets.

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --header-file=HEADER-FILE  The contents of this file will be added to the top of all generated files.
  --filename-managed="zz_generated.managed.go"
                             The filename of generated managed resource files.
  --filename-resolvers="zz_generated.resolvers.go"
                             The filename of generated reference resolver files.
  --filename-managed-list="zz_generated.managedlist.go"
                             The filename of generated managed list resource files.
  --filename-pc="zz_generated.pc.go"
                             The filename of generated provider config files.
  --filename-pcu="zz_generated.pcu.go"
                             The filename of generated provider config usage files.
  --filename-pcu-list="zz_generated.pculist.go"
                             The filename of generated provider config usage files.

Args:
  [<packages>]  Package(s) for which to generate methods, for example github.com/crossplane/crossplane/apis/...
```

[Crossplane]: https://crossplane.io
[`resource.Managed`]: https://godoc.org/github.com/crossplane/crossplane-runtime/v2/pkg/resource#Managed
[`ResourceSpec`]: https://godoc.org/github.com/crossplane/crossplane-runtime/v2/apis/common/v1#ResourceSpec
[`ResourceStatus`]: https://godoc.org/github.com/crossplane/crossplane-runtime/v2/apis/common/v1#ResourceStatus
[Provider Development Guide]: https://github.com/crossplane/crossplane/blob/master/contributing/guide-provider-development.md#defining-resource-kinds