package faceDetection

// package main

import (
	"encoding/json"
	"log"

	pigo "github.com/esimov/pigo/core"
)

type FaceDetect struct {
	classifier    *pigo.Pigo
	cascadeParams *pigo.CascadeParams
	imageParams   *pigo.ImageParams
	Cols          int
	Rows          int
	MinSize       int
	MaxSize       int
	ShiftFactor   float64
	ScaleFactor   float64
	Angle         float64
	IOThreshold   float64
	Faces         string
}

func InitFaceDetect(cas []byte) *FaceDetect {

	var __classifier = pigo.NewPigo()
	var _classifier, error = __classifier.Unpack(cas)

	var _imageParams = &pigo.ImageParams{
		Pixels: []byte(""),
		Rows:   0,
		Cols:   0,
		Dim:    0,
	}

	var _cascadeParams = &pigo.CascadeParams{
		MinSize:     100,
		MaxSize:     400,
		ShiftFactor: 0.1,
		ScaleFactor: 1.0,
		ImageParams: *_imageParams,
	}

	if error != nil {
		log.Fatalf("Error unpacking the cascade file: %s", error)
	}
	return &FaceDetect{
		classifier:    _classifier,
		cascadeParams: _cascadeParams,
		imageParams:   _imageParams,
		Cols:          0,
		Rows:          0,
		MinSize:       100,
		MaxSize:       400,
		ShiftFactor:   0.1,
		ScaleFactor:   1.0,
		Angle:         0.0,
		IOThreshold:   0.0,
		Faces:         "",
	}
}

func (faceDetect *FaceDetect) GetFacesDetect(pixels []byte, cols, rows int) {

	faceDetect.Cols = cols
	faceDetect.Rows = rows
	faceDetect.imageParams = &pigo.ImageParams{
		Pixels: pixels,
		Rows:   faceDetect.Rows,
		Cols:   faceDetect.Cols,
		Dim:    faceDetect.Cols,
	}
	faceDetect.cascadeParams = &pigo.CascadeParams{
		MinSize:     faceDetect.MinSize,
		MaxSize:     faceDetect.MaxSize,
		ShiftFactor: faceDetect.ShiftFactor,
		ScaleFactor: faceDetect.ScaleFactor,
		ImageParams: *faceDetect.imageParams,
	}

	dets := faceDetect.classifier.RunCascade(*faceDetect.cascadeParams, faceDetect.Angle)
	dets = faceDetect.classifier.ClusterDetections(dets, faceDetect.IOThreshold)

	if len(dets) > 0 {
		facesResult, _ := json.Marshal(dets)
		faceDetect.Faces = string(facesResult)
	}
}
