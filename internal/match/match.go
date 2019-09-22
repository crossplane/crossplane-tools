package match

import (
	"go/types"

	"github.com/negz/angryjet/internal/comments"
	"github.com/negz/angryjet/internal/fields"
)

type Object func(o types.Object) bool

func Managed() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta(),
			fields.IsObjectMeta(),
			fields.IsSpec().And(fields.HasFieldThat(fields.IsEmbedded(fields.IsResourceSpec()))),
			fields.IsStatus().And(fields.HasFieldThat(fields.IsEmbedded(fields.IsResourceStatus()))),
		)
	}
}

func NonPortableClass() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta(),
			fields.IsObjectMeta(),
			fields.IsSpecTemplate().And(fields.HasFieldThat(fields.IsEmbedded(fields.IsNonPortableClassSpecTemplate()))),
		)
	}
}

func Claim() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta(),
			fields.IsObjectMeta(),
			fields.IsSpec().And(fields.HasFieldThat(fields.IsEmbedded(fields.IsResourceClaimSpec()))),
			fields.IsResourceClaimStatus(),
		)
	}
}

func PortableClass() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta(),
			fields.IsObjectMeta(),
			fields.IsEmbedded(fields.IsPortableClass()),
		)
	}
}

func HasMethod(name string) Object {
	return func(o types.Object) bool {
		s := types.NewMethodSet(types.NewPointer(o.Type()))
		for i := 0; i < s.Len(); i++ {
			if s.At(i).Obj().Name() == name {
				return true
			}
		}
		return false
	}
}

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

func DoesNotHaveMarker(c comments.Comments, k, v string) Object {
	return func(o types.Object) bool {
		return !HasMarker(c, k, v)(o)
	}
}

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
