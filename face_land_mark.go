package faceDetection

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	pigo "github.com/esimov/pigo/core"
)

type FaceLandMark struct {
	faceClassifier   *pigo.Pigo
	puplocClassifier *pigo.PuplocCascade
	mouthCascade     []string
	eyeCascades      []string
	flpcs            map[string][]*pigo.FlpCascade
	Cols             int
	Rows             int
	MinSize          int
	MaxSize          int
	ShiftFactor      float64
	ScaleFactor      float64
	Angle            float64
	IOThreshold      float64
	Faces            string
	HolesFace        string
}

func InitFaceLandMark(facecascade, puplocCascade, lp38, lp42, lp44, lp46, lp81, lp82, lp84, lp93, lp312 []byte) *FaceLandMark {
	//face classifier
	var __classifier = pigo.NewPigo()
	var _classifier, errorFace = __classifier.Unpack(facecascade)
	if errorFace != nil {
		log.Fatalf("Error unpacking the cascade file: %s", errorFace)
	}
	//puploc classifier
	__puplocClassifier := &pigo.PuplocCascade{}
	_puplocClassifier, errPuploc := __puplocClassifier.UnpackCascade(puplocCascade)
	if errPuploc != nil {
		log.Fatalf("Error unpacking the puploc cascade file: %s", errPuploc)
	}
	//flpcs
	_flpcs := make(map[string][]*pigo.FlpCascade, 9)
	var errFlc error

	////lb38
	__puplocClassifier, errFlc = _puplocClassifier.UnpackCascade(lp38)
	if errFlc != nil {
		log.Fatalf("Error unpacking the puploc cascade file: %s", errFlc)
	}
	_flpcs["lp38"] = append(_flpcs["lp38"], &pigo.FlpCascade{
		PuplocCascade: __puplocClassifier,
	})

	////lp42
	__puplocClassifier, errFlc = _puplocClassifier.UnpackCascade(lp42)
	if errFlc != nil {
		log.Fatalf("Error unpacking the puploc cascade file: %s", errFlc)
	}
	_flpcs["lp42"] = append(_flpcs["lp42"], &pigo.FlpCascade{
		PuplocCascade: __puplocClassifier,
	})

	////lb44
	__puplocClassifier, errFlc = _puplocClassifier.UnpackCascade(lp44)
	if errFlc != nil {
		log.Fatalf("Error unpacking the puploc cascade file: %s", errFlc)
	}
	_flpcs["lp44"] = append(_flpcs["lp44"], &pigo.FlpCascade{
		PuplocCascade: __puplocClassifier,
	})
	////lp46
	__puplocClassifier, errFlc = _puplocClassifier.UnpackCascade(lp46)
	if errFlc != nil {
		log.Fatalf("Error unpacking the puploc cascade file: %s", errFlc)
	}
	_flpcs["lp46"] = append(_flpcs["lp46"], &pigo.FlpCascade{
		PuplocCascade: __puplocClassifier,
	})

	//lp81
	__puplocClassifier, errFlc = _puplocClassifier.UnpackCascade(lp81)
	if errFlc != nil {
		log.Fatalf("Error unpacking the puploc cascade file: %s", errFlc)
	}
	_flpcs["lp81"] = append(_flpcs["lp81"], &pigo.FlpCascade{
		PuplocCascade: __puplocClassifier,
	})
	//lp82
	__puplocClassifier, errFlc = _puplocClassifier.UnpackCascade(lp82)
	if errFlc != nil {
		log.Fatalf("Error unpacking the puploc cascade file: %s", errFlc)
	}
	_flpcs["lp82"] = append(_flpcs["lp82"], &pigo.FlpCascade{
		PuplocCascade: __puplocClassifier,
	})
	//lp84
	__puplocClassifier, errFlc = _puplocClassifier.UnpackCascade(lp84)
	if errFlc != nil {
		log.Fatalf("Error unpacking the puploc cascade file: %s", errFlc)
	}
	_flpcs["lp84"] = append(_flpcs["lp84"], &pigo.FlpCascade{
		PuplocCascade: __puplocClassifier,
	})

	//lp93
	__puplocClassifier, errFlc = _puplocClassifier.UnpackCascade(lp93)
	if errFlc != nil {
		log.Fatalf("Error unpacking the puploc cascade file: %s", errFlc)
	}
	_flpcs["lp93"] = append(_flpcs["lp93"], &pigo.FlpCascade{
		PuplocCascade: __puplocClassifier,
	})
	//lp312
	__puplocClassifier, errFlc = _puplocClassifier.UnpackCascade(lp312)
	if errFlc != nil {
		log.Fatalf("Error unpacking the puploc cascade file: %s", errFlc)
	}
	_flpcs["lp312"] = append(_flpcs["lp312"], &pigo.FlpCascade{
		PuplocCascade: __puplocClassifier,
	})

	_eyeCascades := []string{"lp46", "lp44", "lp42", "lp38", "lp312"}
	_mouthCascade := []string{"lp93", "lp84", "lp82", "lp81"}

	return &FaceLandMark{
		faceClassifier:   _classifier,
		puplocClassifier: _puplocClassifier,
		flpcs:            _flpcs,
		eyeCascades:      _eyeCascades,
		mouthCascade:     _mouthCascade,
		Cols:             0,
		Rows:             0,
		MinSize:          100,
		MaxSize:          400,
		ShiftFactor:      0.1,
		ScaleFactor:      1.0,
		Angle:            0.0,
		IOThreshold:      0.0,
		Faces:            "",
		HolesFace:        "",
	}
}

func (faceLandMark *FaceLandMark) GetFaceLandMark(pixels []uint8, cols int, rows int) {
	faceLandMark.Rows = rows
	faceLandMark.Cols = cols

	imgParams := &pigo.ImageParams{
		Pixels: pixels,
		Rows:   rows,
		Cols:   cols,
		Dim:    cols,
	}
	cParams := pigo.CascadeParams{
		MinSize:     100,
		MaxSize:     400,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: *imgParams,
	}

	filterResult := faceLandMark.faceClassifier.RunCascade(cParams, 0.0)

	fmt.Println("Calculate the intersection over union (IoU) of two clusters.")
	filterResult = faceLandMark.faceClassifier.ClusterDetections(filterResult, 0.0)
	if len(filterResult) > 0 {
		facesResult, _ := json.Marshal(filterResult)
		faceLandMark.Faces = string(facesResult)

		results := []pigo.Detection{filterResult[0]}
		dets := make([][]int, len(results))

		for i := 0; i < len(results); i++ {
			dets[i] = append(dets[i], results[i].Row, results[i].Col, results[i].Scale, int(results[i].Q), 0)
			chLeftEye := make(chan int, 3)
			chLeftRight := make(chan int, 3)
			chPuploc := make(chan pigo.Puploc, 1)
			chLeftEyePuploc := make(chan pigo.Puploc, 1)
			chRightEyePuploc := make(chan pigo.Puploc, 1)
			var wg sync.WaitGroup
			// var puploc pigo.Puploc
			// var leftEyePuploc pigo.Puploc
			// var rightEyePuploc pigo.Puploc
			wg.Add(2)
			// left eye
			go func() {
				puploc := &pigo.Puploc{
					Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
					Col:      results[i].Col - int(0.185*float32(results[i].Scale)),
					Scale:    float32(results[i].Scale) * 0.4,
					Perturbs: 63,
				}
				leftEye := faceLandMark.puplocClassifier.RunDetector(*puploc, *imgParams, 0.0, false)
				if leftEye.Row > 0 && leftEye.Col > 0 {
					dets[i] = append(dets[i], leftEye.Row, leftEye.Col, int(leftEye.Scale), int(results[i].Q), 1)
					chLeftEye <- leftEye.Row
					chLeftEye <- leftEye.Col
					chLeftEye <- int(leftEye.Scale)
					close(chLeftEye)

					chLeftEyePuploc <- *leftEye
					close(chLeftEyePuploc)
					wg.Done()
				}
			}()

			// right eye
			go func() {
				puploc := &pigo.Puploc{
					Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
					Col:      results[i].Col + int(0.185*float32(results[i].Scale)),
					Scale:    float32(results[i].Scale) * 0.4,
					Perturbs: 63,
				}
				rightEye := faceLandMark.puplocClassifier.RunDetector(*puploc, *imgParams, 0.0, false)
				if rightEye.Row > 0 && rightEye.Col > 0 {
					dets[i] = append(dets[i], rightEye.Row, rightEye.Col, int(rightEye.Scale), int(results[i].Q), 1)
				}
				chLeftRight <- rightEye.Row
				chLeftRight <- rightEye.Col
				chLeftRight <- int(rightEye.Scale)
				close(chLeftRight)

				chPuploc <- *puploc
				close(chPuploc)

				chRightEyePuploc <- *rightEye
				close(chRightEyePuploc)
				wg.Done()
			}()

			wg.Wait()
			// puploc = <-chPuploc
			// leftEyePuploc = <-chLeftEyePuploc
			// rightEyePuploc = <-chRightEyePuploc

			// phase 2
			// wg.Add(3)
			// go func() {
			// 	// fmt.Println("-------- eye point-------")
			// 	// fmt.Println(len(faceLandMark.eyeCascades))
			// 	for _, eye := range faceLandMark.eyeCascades {
			// 		for _, flpc := range faceLandMark.flpcs[eye] {
			// 			flp := flpc.GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, false)
			// 			if flp.Row > 0 && flp.Col > 0 {
			// 				dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
			// 			}
			// 			flp = flpc.GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, true)
			// 			if flp.Row > 0 && flp.Col > 0 {
			// 				dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
			// 			}
			// 		}
			// 	}
			// 	wg.Done()
			// }()
			// go func() {
			// 	// fmt.Println("-------- eye point-------")
			// 	// fmt.Println(len(faceLandMark.mouthCascade))
			// 	for _, mouth := range faceLandMark.mouthCascade {
			// 		for _, flpc := range faceLandMark.flpcs[mouth] {
			// 			flp := flpc.GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, false)
			// 			if flp.Row > 0 && flp.Col > 0 {
			// 				dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
			// 			}
			// 		}
			// 	}
			// 	wg.Done()
			// }()
			// go func() {
			// 	flp := faceLandMark.flpcs["lp84"][0].GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, true)
			// 	if flp.Row > 0 && flp.Col > 0 {
			// 		dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
			// 	}
			// 	wg.Done()
			// }()
			// wg.Wait()

		}
		landMarkResult, _ := json.Marshal(dets)
		faceLandMark.HolesFace = string(landMarkResult)
	}
}
