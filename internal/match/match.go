/*
Copyright 2019 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package match identifies Go types as common Crossplane resources.
package match

import (
	"go/types"

	"github.com/crossplane/crossplane-tools/internal/comments"
	"github.com/crossplane/crossplane-tools/internal/fields"
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

// ManagedV2 returns an Object matcher that returns true if the supplied Object is
// a v2-style Crossplane managed resource.
func ManagedV2() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpec().And(fields.HasFieldThat(
				fields.IsResourceV2Spec().And(fields.IsEmbedded()),
			)),
			fields.IsStatus().And(fields.HasFieldThat(
				fields.IsResourceStatus().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ManagedList returns an Object matcher that returns true if the supplied
// Object is a list of Crossplane managed resource.
func ManagedList() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsSpec().And(fields.HasFieldThat(
					fields.IsResourceSpec().And(fields.IsEmbedded()),
				)),
				fields.IsStatus().And(fields.HasFieldThat(
					fields.IsResourceStatus().And(fields.IsEmbedded()),
				)),
			)),
		)
	}
}

// ManagedListV2 returns an Object matcher that returns true if the supplied
// Object is a list of Crossplane managed resource.
func ManagedListV2() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsSpec().And(fields.HasFieldThat(
					fields.IsResourceV2Spec().And(fields.IsEmbedded()),
				)),
				fields.IsStatus().And(fields.HasFieldThat(
					fields.IsResourceStatus().And(fields.IsEmbedded()),
				)),
			)),
		)
	}
}

// ProviderConfig returns an Object matcher that returns true if the supplied
// Object is a Crossplane ProviderConfig.
func ProviderConfig() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpec(),
			fields.IsStatus().And(fields.HasFieldThat(
				fields.IsProviderConfigStatus().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ProviderConfigUsage returns an Object matcher that returns true if the supplied
// Object is a Crossplane ProviderConfigUsage.
func ProviderConfigUsage() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsProviderConfigUsage().And(fields.IsEmbedded()),
		)
	}
}

// ProviderConfigUsageV2 returns an Object matcher that returns true if the supplied
// Object is a Crossplane v2 style ProviderConfigUsage.
func ProviderConfigUsageV2() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsTypedProviderConfigUsage().And(fields.IsEmbedded()),
		)
	}
}

// ProviderConfigUsageList returns an Object matcher that returns true if the
// supplied Object is a list of Crossplane provider config usages.
func ProviderConfigUsageList() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsProviderConfigUsage().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ProviderConfigUsageListV2 returns an Object matcher that returns true if the
// supplied Object is a list of Crossplane v2 style provider config usages.
func ProviderConfigUsageListV2() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsTypedProviderConfigUsage().And(fields.IsEmbedded()),
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
