package processors

// #include <stdlib.h>
import "C"
import (
	"unsafe"
)
func sliceToCArray[T, V any](slice []T, to V) unsafe.Pointer {
	count := C.int(len(slice))
	cArray := C.malloc(C.size_t(count) * C.size_t(unsafe.Sizeof(to)))
	// convert the C array to a Go Array so we can index it
	a := (*[1<<27 - 1]V)(cArray)
	// fill the C array with values
	for index, value := range slice {
		// the following code coerces value to an interface compatible with V
		// where V is a C type compatible with go type T
		v := *((*V)((unsafe.Pointer(&value))))
		if(index == len(slice)) { break }
		a[index] = V(v)
	}

	return cArray
}
