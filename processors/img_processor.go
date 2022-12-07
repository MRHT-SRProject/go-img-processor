package processors

// #cgo CFLAGS: -O2
// #cgo CXXFLAGS: -O2
// #cgo LDFLAGS: -lopencv_core -lopencv_photo -lopencv_imgproc -lopencv_imgcodecs -lopencv_video
// #include "img_processor.h"
// #include <stdlib.h>
import "C"

import (
	"bytes"
	"errors"
	"image"
	"unsafe"

	"github.com/lmittmann/ppm"
	libraw "github.com/richbai90/golibraw"
)

type ColormapTypes uint8
type CMat C.CMat

const (
	COLORMAP_AUTUMN           ColormapTypes = 0
	COLORMAP_BONE             ColormapTypes = 1
	COLORMAP_JET              ColormapTypes = 2
	COLORMAP_WINTER           ColormapTypes = 3
	COLORMAP_RAINBOW          ColormapTypes = 4
	COLORMAP_OCEAN            ColormapTypes = 5
	COLORMAP_SUMMER           ColormapTypes = 6
	COLORMAP_SPRING           ColormapTypes = 7
	COLORMAP_COOL             ColormapTypes = 8
	COLORMAP_HSV              ColormapTypes = 9
	COLORMAP_PINK             ColormapTypes = 10
	COLORMAP_HOT              ColormapTypes = 11
	COLORMAP_PARULA           ColormapTypes = 12
	COLORMAP_MAGMA            ColormapTypes = 13
	COLORMAP_INFERNO          ColormapTypes = 14
	COLORMAP_PLASMA           ColormapTypes = 15
	COLORMAP_VIRIDIS          ColormapTypes = 16
	COLORMAP_CIVIDIS          ColormapTypes = 17
	COLORMAP_TWILIGHT         ColormapTypes = 18
	COLORMAP_TWILIGHT_SHIFTED ColormapTypes = 19
	COLORMAP_TURBO            ColormapTypes = 20
	COLORMAP_DEEPGREEN        ColormapTypes = 21
)

func GetImageFromRaw(rawImg []byte) (image.Image, libraw.ImgMetadata, error) {
	return libraw.RawBuffer2Image(rawImg)
}

func StackImages(imgs ...image.Image) (C.CMat, error) {

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

	stacked := C.stackImages(cimgs, C.size_t(cImgLen))

	defer func() {
		for _, img := range cImages {
			C.free(unsafe.Pointer(img.pixels))
		}

		C.free(unsafe.Pointer(cimgs))
	}()

	err := errFromResp(stacked)
	if err != nil {
		return C.CMat{}, err
	}

	return C.CMat{
		mat: stacked.data,
	}, nil
}

func Colorize(img C.CMat, colormap ColormapTypes) (C.CMat, error) {
	resp := C.colorize(img, C.uchar(colormap))
	err := errFromResp(resp)
	if err != nil {
		return C.CMat{}, err
	}

	return C.CMat{
		mat: resp.data,
	}, nil
}

func GrayScale(imgs ...image.Image) ([]image.Image, error) {
	gsimgs := make([]image.Image, len(imgs))

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

		resp := C.grayscale(cImage)
		err := errFromResp(resp)
		if err != nil {
			return imgs, err
		}
		gs := C.CMat{mat: resp.data}
		cpixels := C.getPixels(gs)
		gsimg := image.NewGray(rect)
		gsimg.Pix = cArrToSlice(unsafe.Pointer(cpixels.pixels), uint8(0), uint(cpixels.len))
		gsimgs[i] = gsimg

		defer func() {
			C.free(gs.mat)
		}()
	}

	return gsimgs, nil
}

func CMatToImg(mat C.CMat) (image.Image, error) {
	buf := C.cMatToImg(mat)
	defer C.free(unsafe.Pointer(buf.data))
	err := errFromResp(buf)
	if err != nil {
		return nil, err
	}
	s := cArrToSlice(unsafe.Pointer(buf.data), uint8(0), uint(buf.len))
	ppmData := bytes.NewBuffer(s)
	return ppm.Decode(ppmData)
}

func errFromResp(resp C.Response) error {
	if resp.status != 0 {
		str := cArrToSlice(resp.data, ' ', uint(resp.len))
		return errors.New(string(str))
	}

	return nil
}
