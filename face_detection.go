package faceDetection

// package main

import (
	"encoding/json"
	"fmt"
	"log"

	pigo "github.com/esimov/pigo/core"
)

// type FacePoint struct {
// 	row   int
// 	col   int
// 	scale int
// 	q     int
// }

//-----------------------------------------------------------------------------
type MobilePigo struct {
	classifier    *pigo.Pigo
	cascadwParams *pigo.CascadeParams
	Cols          int
	Rows          int
	Faces         string
	NumFace       int
}

func InitFaceDetect(cas []byte) *MobilePigo {

	var __classifier = pigo.NewPigo()
	var _classifier, error = __classifier.Unpack(cas)
	if error != nil {
		log.Fatalf("Error unpacking the cascade file: %s", error)
	}
	var _cascadwParams = pigo.CascadeParams{}
	return &MobilePigo{
		classifier:    _classifier,
		cascadwParams: &_cascadwParams,
		Cols:          0,
		Rows:          0,
		Faces:         "",
		NumFace:       0,
	}
}

func (mobilego *MobilePigo) GetFacesDetect(pixels []byte, cols, rows int) {

	// var img, err = jpeg.Decode(bytes.NewReader(bytesImage))
	// if err != nil {
	// 	log.Fatalf("Error reading the cascade file: %s", err)
	// }
	// fmt.Println("start")

	// var pixels = pigo.RgbToGrayscale(img)
	fmt.Println(len(pixels))

	// cols, rows := img.Bounds().Max.X, img.Bounds().Max.Y
	mobilego.Cols = cols
	mobilego.Rows = rows

	mobilego.cascadwParams.MinSize = 100
	mobilego.cascadwParams.MaxSize = 600
	mobilego.cascadwParams.ShiftFactor = 0.15
	mobilego.cascadwParams.ScaleFactor = 1.1
	mobilego.cascadwParams.ImageParams = pigo.ImageParams{
		Pixels: pixels,
		Rows:   rows,
		Cols:   cols,
		Dim:    cols,
	}

	dets := mobilego.classifier.RunCascade(*mobilego.cascadwParams, 0.0)
	dets = mobilego.classifier.ClusterDetections(dets, 0.0)

	if len(dets) > 0 {
		if mobilego.NumFace > len(dets) {
			mobilego.NumFace = len(dets)
		}
		// var faces []FacePoint
		// for i := 0; i < mobilego.NumFace; i++ {
		// 	faces = append(faces, FacePoint{
		// 		row:   dets[i].Row,
		// 		col:   dets[i].Col,
		// 		q:     int(dets[i].Q),
		// 		scale: dets[i].Scale,
		// 	})
		// }
		facesResult, _ := json.Marshal(dets)
		fmt.Println(string(facesResult))
		mobilego.Faces = string(facesResult)
	}
	// else {
	// mobilego.Faces = ""
	// fmt.Println("No face founded")
	// }
}

//-----------------------------------------------------------------------------
