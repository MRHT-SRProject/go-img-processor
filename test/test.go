package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/MRHT-SRProject/go-img-processor/processors"
)

func main() {
	csv_file, err := os.Open("/home/rich/code/camcapture/test/files.csv")
	handleError(err, "Error opening csv file. ")
	csv_reader := csv.NewReader(csv_file)
	file_list, err := csv_reader.ReadAll()
	handleError(err, "Failed to parse csv file")
	images := make([]image.RGBA, len(file_list))
	for i, info := range file_list {
		fname := info[0]
		raw, err := os.Open("/home/rich/code/camcapture/test/" + fname)
		handleError(err, "Failed to open ", fname)
		finfo, err := raw.Stat()
		handleError(err, "Failed to stat ", fname)
		buffer := make([]byte, finfo.Size())
		_, err = raw.Read(buffer)
		handleError(err, "Failed to read ", fname)
		raw.Close()
		img, _, err := processors.GetImageFromRaw(buffer)
		handleError(err, "Failed to convert raw file ", fname, " to image")
		rgba, _ := img.(*image.RGBA)
		images[i] = *rgba;
	}

	gsimgs := processors.GrayScale(images...)
	for i, v := range gsimgs {
		f, err := os.Create(fmt.Sprintf("%d.png", i));
		handleError(err, "Failed to create ", i, ".png")
		png.Encode(f, &v);
	}
	print(gsimgs)
}

func handleError(err error, v ...any) {
	if err != nil {
		v = append(v, "ERROR: ", err.Error())
		log.Fatal(v...)
	}
}