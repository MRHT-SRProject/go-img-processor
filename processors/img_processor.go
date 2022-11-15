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

type Image struct {
	image.Image
	Exposure float32
}

func GetImageFromRaw(rawImg []byte) (image.Image, libraw.ImgMetadata, error) {
	return libraw.RawBuffer2Image(rawImg)
}

func StackImages(imgs ...image.Image) {
	cImgLen := len(imgs)
	cImages := make([]C.Image, cImgLen)
	for i, iimg := range imgs {
		img := iimg.(*image.RGBA)
		imgDataC := sliceToCArray(img.Pix, C.uchar(0))
		cImage := C.Image{
			pixels: (*C.uchar)(imgDataC),
			x0:     C.uint(img.Rect.Min.X),
			x1:     C.uint(img.Rect.Max.X),
			y0:     C.uint(img.Rect.Min.Y),
			y1:     C.uint(img.Rect.Max.Y),
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

func GrayScale(imgs ...image.RGBA) []image.Gray {
	gsimgs := make([]image.Gray, len(imgs))

	for i, img := range imgs {
		rect := img.Rect
		imgDataC := sliceToCArray(img.Pix, C.uchar(0))
		cImage := C.Image{
			pixels: (*C.uchar)(imgDataC),
			x0:     C.uint(img.Rect.Min.X),
			x1:     C.uint(img.Rect.Max.X),
			y0:     C.uint(img.Rect.Min.Y),
			y1:     C.uint(img.Rect.Max.Y),
		}

		gs := C.grayscale(cImage)
		cpixels := C.getPixels(gs)
		gsimg := image.NewGray(rect)
		gsimg.Pix = unsafe.Slice((*uint8)(cpixels.pixels), cpixels.len);
		gsimgs[i] = *gsimg
		
	}

	return gsimgs
}
