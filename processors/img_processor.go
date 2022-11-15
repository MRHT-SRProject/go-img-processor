package processors

// #cgo CFLAGS: -O2
// #cgo CXXFLAGS: -O2
// #cgo LDFLAGS: -lopencv_core -lopencv_photo -lopencv_imgproc -lopencv_imgcodecs -lopencv_video
// #include "img_processor.h"
// #include <stdlib.h>
import "C"

import (
	"image"
	"unsafe"

	libraw "github.com/richbai90/golibraw"
)

func GetImageFromRaw(rawImg []byte) (image.Image, libraw.ImgMetadata, error) {
	return libraw.RawBuffer2Image(rawImg)
}

func StackImages(imgs ...image.Image) {

	cImgLen := len(imgs)
	cImages := make([]C.Image, cImgLen)
	for i, iimg := range imgs {
		img := NewGenericImage(iimg)
		pix, ok := img.GetPix()
		if !ok {
			continue
		}
		rect, ok := img.GetRect()
		if !ok {
			continue
		}
		imgDataC := sliceToCArray(pix, C.uchar(0))
		cImage := C.Image{
			pixels: (*C.uchar)(imgDataC),
			x0:     C.uint(rect.Min.X),
			x1:     C.uint(rect.Max.X),
			y0:     C.uint(rect.Min.Y),
			y1:     C.uint(rect.Max.Y),
		}
		cImages[i] = cImage
	}
	cimgs := (*C.Image)(sliceToCArray(cImages, C.Image{}))
	C.stackImages(cimgs, C.size_t(cImgLen))
	defer func() {
		for _, img := range cImages {
			C.free(unsafe.Pointer(img.pixels))
		}

		C.free(unsafe.Pointer(cimgs))
	}()

}

func GrayScale(imgs ...image.Image) []image.Gray {
	gsimgs := make([]image.Gray, len(imgs))

	for i, iimg := range imgs {
		img := NewGenericImage(iimg)
		rect, _ := img.GetRect()
		pix, _ := img.GetPix()
		imgDataC := sliceToCArray(pix, C.uchar(0))
		cImage := C.Image{
			pixels: (*C.uchar)(imgDataC),
			x0:     C.uint(rect.Min.X),
			x1:     C.uint(rect.Max.X),
			y0:     C.uint(rect.Min.Y),
			y1:     C.uint(rect.Max.Y),
		}

		gs := C.grayscale(cImage)
		cpixels := C.getPixels(gs)
		gsimg := image.NewGray(rect)
		gsimg.Pix = unsafe.Slice((*uint8)(cpixels.pixels), cpixels.len)
		gsimgs[i] = *gsimg

	}

	return gsimgs
}
