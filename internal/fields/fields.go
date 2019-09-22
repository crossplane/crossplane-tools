package fields

import (
	"go/types"
	"strings"
)

// Field names.
const (
	NameTypeMeta                     = "TypeMeta"
	NameObjectMeta                   = "ObjectMeta"
	NameSpec                         = "Spec"
	NameSpecTemplate                 = "SpecTemplate"
	NameStatus                       = "Status"
	NameResourceSpec                 = "ResourceSpec"
	NameResourceStatus               = "ResourceStatus"
	NameResourceClaimSpec            = "ResourceClaimSpec"
	NameNonPortableClassSpecTemplate = "NonPortableClassSpecTemplate"
	NamePortableClass                = "PortableClass"
)

// Field type suffixes.
const (
	TypeSuffixTypeMeta                     = "k8s.io/apimachinery/pkg/apis/meta/v1.TypeMeta"
	TypeSuffixObjectMeta                   = "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"
	TypeSuffixSpec                         = NameSpec
	TypeSuffixSpecTemplate                 = NameSpecTemplate
	TypeSuffixStatus                       = NameStatus
	TypeSuffixResourceSpec                 = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.ResourceSpec"
	TypeSuffixResourceStatus               = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.ResourceStatus"
	TypeSuffixResourceClaimSpec            = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.ResourceClaimSpec"
	TypeSuffixResourceClaimStatus          = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.ResourceClaimStatus"
	TypeSuffixNonPortableClassSpecTemplate = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.NonPortableClassSpecTemplate"
	TypeSuffixPortableClass                = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1.PortableClass"
)

func matches(s *types.Struct, m Matcher) bool {
	for i := 0; i < s.NumFields(); i++ {
		if m(s.Field(i)) {
			return true
		}
	}
	return false
}

func Has(o types.Object, m ...Matcher) bool {
	s, ok := o.Type().Underlying().(*types.Struct)
	if !ok {
		return false
	}
	for _, matcher := range m {
		if !matches(s, matcher) {
			return false
		}
	}
	return true
}

func Get(o types.Object, m Matcher) *types.Var {
	s, ok := o.Type().Underlying().(*types.Struct)
	if !ok {
		return nil
	}
	for i := 0; i < s.NumFields(); i++ {
		if m(s.Field(i)) {
			return s.Field(i)
		}
	}
	return nil
}

type Matcher func(f *types.Var) bool

func (o Matcher) And(m Matcher) Matcher {
	return func(f *types.Var) bool {
		return o(f) && m(f)
	}
}

func IsEmbedded(m Matcher) Matcher {
	return func(f *types.Var) bool {
		if !f.Embedded() {
			return false
		}
		return m(f)
	}
}

func IsNamed(name string) Matcher {
	return func(f *types.Var) bool {
		if !f.IsField() {
			return false
		}
		return f.Name() == name
	}
}

func IsTypeNamed(typeName, name string) Matcher {
	return func(f *types.Var) bool {
		if !IsNamed(name)(f) {
			return false
		}
		return strings.HasSuffix(f.Type().String(), typeName)
	}
}

func IsTypeMeta() Matcher            { return IsTypeNamed(TypeSuffixTypeMeta, NameTypeMeta) }
func IsObjectMeta() Matcher          { return IsTypeNamed(TypeSuffixObjectMeta, NameObjectMeta) }
func IsSpec() Matcher                { return IsTypeNamed(NameSpec, TypeSuffixSpec) }
func IsSpecTemplate() Matcher        { return IsTypeNamed(NameSpecTemplate, TypeSuffixSpecTemplate) }
func IsStatus() Matcher              { return IsTypeNamed(NameStatus, TypeSuffixStatus) }
func IsResourceStatus() Matcher      { return IsTypeNamed(TypeSuffixResourceStatus, NameResourceStatus) }
func IsResourceSpec() Matcher        { return IsTypeNamed(TypeSuffixResourceSpec, NameResourceSpec) }
func IsResourceClaimStatus() Matcher { return IsTypeNamed(TypeSuffixResourceClaimStatus, NameStatus) }
func IsResourceClaimSpec() Matcher {
	return IsTypeNamed(TypeSuffixResourceClaimSpec, NameResourceClaimSpec)
}
func IsNonPortableClassSpecTemplate() Matcher {
	return IsTypeNamed(TypeSuffixNonPortableClassSpecTemplate, NameNonPortableClassSpecTemplate)
}
func IsPortableClass() Matcher {
	return IsTypeNamed(TypeSuffixPortableClass, NamePortableClass)
}

func HasFieldThat(m ...Matcher) Matcher {
	return func(f *types.Var) bool {
		return Has(f, m...)
	}
}
