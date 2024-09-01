package wca

import (
	"errors"

	"github.com/go-ole/go-ole"
)

type PROPVARIANT struct {
	ole.VARIANT
}

func (v PROPVARIANT) String() string {
	return pvString(v.Val)
}

func (v PROPVARIANT) Bool() (bool, error) {
	a := v.Value()
	if b, ok := a.(bool); ok {
		return b, nil
	}
	return false, errors.New("variant is not bool")
}

func NewStringPropVariant(s string) (PROPVARIANT, error) {
	return stringToPropVariant(s)
}

func NewBoolPropVariant(b bool) PROPVARIANT {
	v := ole.VARIANT{}
	v.VT = ole.VT_BOOL
	switch b {
	case true:
		v.Val = 0xffff
	case false:
		v.Val = 0
	}
	return PROPVARIANT{v}
}

// NewEmptyPropVariant creates a new empty PROPVARIANT. If we want to remove a
// property from a property store, we do it by setting the property to an empty
// PROPVARIANT.
func NewEmptyPropVariant() PROPVARIANT {
	v := ole.VARIANT{}
	v.VT = ole.VT_EMPTY
	return PROPVARIANT{v}
}
