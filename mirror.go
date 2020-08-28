package mirror

import (
	"reflect"
	"unsafe"
)

// DeepCopy implements a deep copy algorithm for Go objects using the reflect package.
// Use DeepCopyInto to initialize this struct.
type DeepCopy struct {
	target           reflect.Value
	ignoreZeroValues bool
	ignoreUnexported bool
}

// DeepCopyInto creates a new DeepCopy configured to copy into the target object.
// The target object must be a pointer, otherwise it cannot be mutated. If it is not a pointer, this function will
// panic. By default, the DeepCopy returned by this function is set to not ignore zero values and to not ignored
// unexported  fields.
func DeepCopyInto(target interface{}) *DeepCopy {
	tval := reflect.ValueOf(target)
	if tval.Kind() != reflect.Ptr {
		panic("cannot mutate target")
	}

	return &DeepCopy{
		target:           tval,
		ignoreZeroValues: false,
		ignoreUnexported: false,
	}
}

// SetIgnoreZeroValues sets whether or not this DeepCopy object should ignore zero values.
func (d *DeepCopy) SetIgnoreZeroValues(ignoreZeroValues bool) *DeepCopy {
	d.ignoreZeroValues = ignoreZeroValues
	return d
}

// SetIgnoreUnexported sets whether or not this DeepCopy object should ignore unexported fields.
func (d *DeepCopy) SetIgnoreUnexported(ignoreUnexported bool) *DeepCopy {
	d.ignoreUnexported = ignoreUnexported
	return d
}

// From performs a deep copy of the src object into the target object.
// The src object must be of the same type as the target object. It does not have to be a pointer but it
// can be. The src object will not be mutated.
func (d *DeepCopy) From(src interface{}) {
	sval := reflect.ValueOf(src)

	// Cast T to *T
	if sval.Kind() != reflect.Ptr {
		temp := reflect.New(reflect.TypeOf(src))
		temp.Elem().Set(sval)
		sval = temp
	}
	if d.target.Type() != sval.Type() {
		panic("different types between target and source")
	}

	d.performDeepCopy(d.target, sval)
}

func (d *DeepCopy) performDeepCopy(target reflect.Value, src reflect.Value) {
	switch src.Kind() {
	case reflect.Invalid:
	case reflect.Array:
		length := src.Len()
		for i := 0; i < length; i++ {
			d.performDeepCopy(target.Index(i), src.Index(i))
		}
	case reflect.Slice:
		length := src.Len()
		newSlice := reflect.MakeSlice(src.Type(), length, src.Cap())
		for i := 0; i < length; i++ {
			d.performDeepCopy(newSlice.Index(i), src.Index(i))
		}
		target.Set(newSlice)
	case reflect.Map:
		newMap := reflect.MakeMapWithSize(src.Type(), src.Len())
		iter := src.MapRange()
		for iter.Next() {
			newMap.SetMapIndex(iter.Key(), iter.Value())
		}
		target.Set(newMap)
	case reflect.Struct:
		length := src.NumField()
		for i := 0; i < length; i++ {
			tfield := target.Field(i)
			sfield := src.Field(i)

			isUnexported := sfield.CanAddr() && !sfield.CanSet()
			if isUnexported && d.ignoreUnexported {
				continue
			}

			if isUnexported {
				tfield = exportUnexportedField(tfield)
				sfield = exportUnexportedField(sfield)
			}

			d.performDeepCopy(tfield, sfield)
		}
	case reflect.Ptr:
		if target.IsNil() && !src.IsNil() {
			// If target is a nil pointer and src is not nil, create a zeroed object for target.
			target.Set(reflect.New(target.Type().Elem()))
		}
		d.performDeepCopy(target.Elem(), src.Elem())
	default:
		if target.CanSet() && !(d.ignoreZeroValues && src.IsZero()) {
			target.Set(src)
		}
	}
}

func exportUnexportedField(field reflect.Value) reflect.Value {
	// Trick to mutate unexported field: https://stackoverflow.com/a/43918797
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}
