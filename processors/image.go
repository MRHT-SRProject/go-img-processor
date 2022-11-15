package processors

import "image"

// Generic image wraps the image.Image interface
// And provides methods for getting the Pix and Rect
// Values without typecasting
type GenericImage struct {
	image.Image
}

// Get the pixel vaues from the image
func (img GenericImage) GetPix() ([]uint8, bool) {
	switch img.Image.(type) {
	case *image.Alpha:
		return img.Image.(*image.Alpha).Pix, true
	case *image.CMYK:
		return img.Image.(*image.CMYK).Pix, true
	case *image.Gray:
		return img.Image.(*image.Gray).Pix, true
	case *image.NRGBA:
		return img.Image.(*image.NRGBA).Pix, true
	case *image.RGBA:
		return img.Image.(*image.RGBA).Pix, true
	default:
		return nil, false
	}
}

// Get the rectangle of the image
func (img GenericImage) GetRect() (image.Rectangle, bool) {
	switch img.Image.(type) {
	case *image.Alpha:
		return img.Image.(*image.Alpha).Rect, true
	case *image.CMYK:
		return img.Image.(*image.CMYK).Rect, true
	case *image.Gray:
		return img.Image.(*image.Gray).Rect, true
	case *image.NRGBA:
		return img.Image.(*image.NRGBA).Rect, true
	case *image.RGBA:
		return img.Image.(*image.RGBA).Rect, true
	default:
		return image.Rectangle{}, false
	}
}

func NewGenericImage(img image.Image) GenericImage {
	return GenericImage{Image: img}
}
