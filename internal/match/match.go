// Package match identifies Go types as common Crossplane resources.
package match

import (
	"go/types"

	"github.com/negz/angryjet/internal/comments"
	"github.com/negz/angryjet/internal/fields"
)

// An Object matcher is a function that returns true if the supplied object
// matches.
type Object func(o types.Object) bool

// Managed returns an Object matcher that returns true if the supplied Object is
// a Crossplane managed resource.
func Managed() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpec().And(fields.HasFieldThat(
				fields.IsResourceSpec().And(fields.IsEmbedded()),
			)),
			fields.IsStatus().And(fields.HasFieldThat(
				fields.IsResourceStatus().And(fields.IsEmbedded()),
			)),
		)
	}
}

// NonPortableClass returns an Object matcher that returns true if the supplied
// Object is a Crossplane non-portable resource class.
func NonPortableClass() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpecTemplate().And(fields.HasFieldThat(
				fields.IsNonPortableClassSpecTemplate().And(fields.IsEmbedded()),
			)),
		)
	}
}

// Claim returns an Object matcher that returns true if the supplied Object is a
// Crossplane resource claim.
func Claim() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpec().And(fields.HasFieldThat(
				fields.IsResourceClaimSpec().And(fields.IsEmbedded()),
			)),
			fields.IsResourceClaimStatus(),
		)
	}
}

// PortableClass returns an Object matcher that returns true if the supplied
// Object is a Crossplane portable resource class.
func PortableClass() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsPortableClass().And(fields.IsEmbedded()),
		)
	}
}

// PortableClassList returns an Object matcher that returns true if the supplied
// Object is a Crossplane portable resource class list.
func PortableClassList() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsListMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsPortableClass().And(fields.IsEmbedded()),
			)),
		)
	}
}

// HasMarker returns an Object matcher that returns true if the supplied Object
// has a comment marker k with the value v. Comment markers are read from the
// supplied Comments.
func HasMarker(c comments.Comments, k, v string) Object {
	return func(o types.Object) bool {
		for _, val := range comments.ParseMarkers(c.For(o))[k] {
			if val == v {
				return true
			}
		}

		for _, val := range comments.ParseMarkers(c.Before(o))[k] {
			if val == v {
				return true
			}
		}

		return false
	}
}

// DoesNotHaveMarker returns and Object matcher that returns true if the
// supplied Object does not have a comment marker k with the value v. Comment
// marker are read from the supplied Comments.
func DoesNotHaveMarker(c comments.Comments, k, v string) Object {
	return func(o types.Object) bool {
		return !HasMarker(c, k, v)(o)
	}
}

// AllOf returns an Object matcher that returns true if all of the supplied
// Object matchers return true.
func AllOf(match ...Object) Object {
	return func(o types.Object) bool {
		for _, fn := range match {
			if !fn(o) {
				return false
			}
		}
		return true
	}
}

// AnyOf returns an Object matcher that returns true if any of the supplied
// Object matchers return true.
func AnyOf(match ...Object) Object {
	return func(o types.Object) bool {
		for _, fn := range match {
			if fn(o) {
				return true
			}
		}
		return false
	}
}
