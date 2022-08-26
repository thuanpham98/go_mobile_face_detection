package faceDetection

import (
	"encoding/json"
	"log"
	"math"

	pigo "github.com/esimov/pigo/core"
)

type point struct {
	x, y int
}

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
	NumFace          int
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
		MinSize:          200,
		MaxSize:          640,
		ShiftFactor:      0.1,
		ScaleFactor:      1.0,
		Angle:            0.0,
		IOThreshold:      0.0,
		NumFace:          0,
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
		MinSize:     200,
		MaxSize:     720,
		ShiftFactor: 0.1,
		ScaleFactor: 1,
		ImageParams: *imgParams,
	}

	results := faceLandMark.faceClassifier.RunCascade(cParams, 0.0)
	results = faceLandMark.faceClassifier.ClusterDetections(results, 0.15)

	if len(results) > 0 {
		facesResult, _ := json.Marshal(results)
		faceLandMark.Faces = string(facesResult)

		// results := filterResult
		dets := make([][]int, len(results))
		var detsRet [][]int
		for i := 0; i < len(results); i++ {

			// dets[i] = append(dets[i], results[i].Row, results[i].Col, results[i].Scale, int(results[i].Q), 0)

			// chPuploc := make(chan pigo.Puploc, 1)
			// chLeftEyePuploc := make(chan pigo.Puploc, 1)
			// chRightEyePuploc := make(chan pigo.Puploc, 1)
			// var wg sync.WaitGroup
			// var puploc pigo.Puploc
			// var leftEyePuploc pigo.Puploc
			// var rightEyePuploc pigo.Puploc
			// wg.Add(1)

			// go func() {

			// left eye
			puploc := pigo.Puploc{
				Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
				Col:      results[i].Col - int(0.185*float32(results[i].Scale)),
				Scale:    float32(results[i].Scale) * 0.4,
				Perturbs: 63,
			}
			leftEye := faceLandMark.puplocClassifier.RunDetector(puploc, *imgParams, 0.0, false)
			// if leftEye.Row > 0 && leftEye.Col > 0 {
			// 	dets[i] = append(dets[i], leftEye.Row, leftEye.Col, int(leftEye.Scale), int(results[i].Q), 1)
			// }
			p1 := &point{x: leftEye.Row, y: leftEye.Col}

			// chLeftEyePuploc <- *leftEye
			// close(chLeftEyePuploc)

			// right eye
			puploc = pigo.Puploc{
				Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
				Col:      results[i].Col + int(0.185*float32(results[i].Scale)),
				Scale:    float32(results[i].Scale) * 0.4,
				Perturbs: 63,
			}
			rightEye := faceLandMark.puplocClassifier.RunDetector(puploc, *imgParams, 0.0, false)
			// if rightEye.Row > 0 && rightEye.Col > 0 {
			// 	dets[i] = append(dets[i], rightEye.Row, rightEye.Col, int(rightEye.Scale), int(results[i].Q), 1)
			// }
			p2 := &point{x: rightEye.Row, y: rightEye.Col}

			// Calculate the lean angle between the pupils.
			angle := math.Atan2(float64(p2.y-p1.y), float64(p2.x-p1.x)) * 180 / math.Pi
			// face
			dets[i] = append(dets[i], results[i].Row, results[i].Col, results[i].Scale, int(results[i].Q), int(angle))

			// 	chRightEyePuploc <- *rightEye
			// 	close(chRightEyePuploc)

			// 	chPuploc <- *puploc
			// 	close(chPuploc)
			// 	wg.Done()
			// }()

			// wg.Wait()

			// puploc = <-chPuploc
			// leftEyePuploc = <-chLeftEyePuploc
			// rightEyePuploc = <-chRightEyePuploc

			// phase 2 --------------------------------------
			flp := faceLandMark.flpcs["lp93"][0].GetLandmarkPoint(leftEye, rightEye, *imgParams, puploc.Perturbs, true)
			if flp.Row > 0 && flp.Col > 0 {
				dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
			}
			// wg.Add(1)
			// go func() {
			// for _, eye := range faceLandMark.eyeCascades {
			// 	for _, flpc := range faceLandMark.flpcs[eye] {
			// 		flp := flpc.GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, false)
			// 		if flp.Row > 0 && flp.Col > 0 {
			// 			dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
			// 		}
			// 		flp = flpc.GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, true)
			// 		if flp.Row > 0 && flp.Col > 0 {
			// 			dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
			// 		}
			// 	}
			// }
			// for _, mouth := range faceLandMark.mouthCascade {
			// 	for _, flpc := range faceLandMark.flpcs[mouth] {
			// 		flp := flpc.GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, false)
			// 		if flp.Row > 0 && flp.Col > 0 {
			// 			dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
			// 		}
			// 	}
			// }
			// 	flp := faceLandMark.flpcs["lp93"][0].GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, true)
			// 	if flp.Row > 0 && flp.Col > 0 {
			// 		dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
			// 	}
			// 	wg.Done()
			// }()
			// wg.Wait()
			if len(dets[i]) == 10 {
				detsRet = append(detsRet, dets[i])
			}
		}
		if len(detsRet) == 0 {
			faceLandMark.Faces = ""
			faceLandMark.HolesFace = ""
			faceLandMark.NumFace = 0

		} else {
			landMarkResult, _ := json.Marshal(detsRet)
			faceLandMark.HolesFace = string(landMarkResult)
			faceLandMark.NumFace = len(detsRet)
		}

	} else {
		faceLandMark.Faces = ""
		faceLandMark.HolesFace = ""
		faceLandMark.NumFace = 0
	}
}
