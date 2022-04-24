package faceDetection

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"sync"

	pigo "github.com/esimov/pigo/core"
)

type MobilePigoLandMark struct {
	faceClassifier   *pigo.Pigo
	puplocClassifier *pigo.PuplocCascade
	flpcs            map[string][]*pigo.FlpCascade
	eyeCascades      []string
	mouthCascade     []string
	Col              int
	Row              int
	Q                int
	Cols             int
	Rows             int
	Scale            int
	eyeLeftcol       int
	eyeLeftRow       int
	eyeRightcol      int
	eyeRightRow      int
	NoseCol          int
	NoseRow          int
}

func InitFaceLandMark(facecascade, puplocCascade, lp38, lp42, lp44, lp46, lp81, lp82, lp84, lp93, lp312 []byte) *MobilePigoLandMark {
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

	return &MobilePigoLandMark{
		eyeCascades:      _eyeCascades,
		mouthCascade:     _mouthCascade,
		faceClassifier:   _classifier,
		puplocClassifier: _puplocClassifier,
		flpcs:            _flpcs,
		Col:              0,
		Row:              0,
		Q:                0,
		Cols:             0,
		Rows:             0,
		Scale:            0,
		eyeLeftcol:       0,
		eyeLeftRow:       0,
		eyeRightcol:      0,
		eyeRightRow:      0,
		NoseCol:          0,
		NoseRow:          0,
	}
}

// func (mobiLM *MobilePigoLandMark) GetFaceLandMark(bytesImage []uint8) {

// 	//// ----------------------------------------processing face Image
// 	var img, err = jpeg.Decode(bytes.NewReader(bytesImage))
// 	if err != nil {
// 		log.Fatalf("Error reading the cascade file: %s", err)
// 	}
// 	fmt.Println("start")

// 	var pixels = pigo.RgbToGrayscale(img)
// 	fmt.Println(len(pixels))
// 	cols, rows := img.Bounds().Max.X, img.Bounds().Max.Y

// 	imgParams := &pigo.ImageParams{
// 		Pixels: pixels,
// 		Rows:   rows,
// 		Cols:   cols,
// 		Dim:    cols,
// 	}
// 	cParams := pigo.CascadeParams{
// 		MinSize:     60,
// 		MaxSize:     600,
// 		ShiftFactor: 0.1,
// 		ScaleFactor: 1.1,
// 		ImageParams: *imgParams,
// 	}
// 	//// ----------------------------------------processing face Detect

// 	filterResult := mobiLM.faceClassifier.RunCascade(cParams, 0.0)

// 	// Calculate the intersection over union (IoU) of two clusters.
// 	filterResult = mobiLM.faceClassifier.ClusterDetections(filterResult, 0.0)
// 	if len(filterResult) > 0 {
// 		results := []pigo.Detection{filterResult[0]}

// 		mobiLM.Scale = filterResult[0].Scale
// 		mobiLM.Row = filterResult[0].Row
// 		mobiLM.Col = filterResult[0].Col
// 		mobiLM.Rows = rows
// 		mobiLM.Cols = cols
// 		mobiLM.Q = int(filterResult[0].Q)

// 		dets := make([][]int, len(results))

// 		for i := 0; i < len(results); i++ {
// 			dets[i] = append(dets[i], results[i].Row, results[i].Col, results[i].Scale, int(results[i].Q), 0)
// 			// left eye
// 			puploc := &pigo.Puploc{
// 				Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
// 				Col:      results[i].Col - int(0.185*float32(results[i].Scale)),
// 				Scale:    float32(results[i].Scale) * 0.4,
// 				Perturbs: 63,
// 			}
// 			leftEye := mobiLM.puplocClassifier.RunDetector(*puploc, *imgParams, 0.0, false)
// 			if leftEye.Row > 0 && leftEye.Col > 0 {
// 				dets[i] = append(dets[i], leftEye.Row, leftEye.Col, int(leftEye.Scale), int(results[i].Q), 1)
// 				mobiLM.eyeLeftRow = leftEye.Row
// 				mobiLM.eyeLeftcol = leftEye.Col
// 			}

// 			// right eye
// 			puploc = &pigo.Puploc{
// 				Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
// 				Col:      results[i].Col + int(0.185*float32(results[i].Scale)),
// 				Scale:    float32(results[i].Scale) * 0.4,
// 				Perturbs: 63,
// 			}

// 			rightEye := mobiLM.puplocClassifier.RunDetector(*puploc, *imgParams, 0.0, false)
// 			if rightEye.Row > 0 && rightEye.Col > 0 {
// 				dets[i] = append(dets[i], rightEye.Row, rightEye.Col, int(rightEye.Scale), int(results[i].Q), 1)
// 				mobiLM.eyeRightRow = rightEye.Row
// 				mobiLM.eyeRightcol = rightEye.Col
// 			}

// 			// Traverse all the eye cascades and run the detector on each of them.
// 			for _, eye := range mobiLM.eyeCascades {
// 				fmt.Println("-------- eye point-------")
// 				fmt.Println(len(mobiLM.flpcs[eye]))
// 				for _, flpc := range mobiLM.flpcs[eye] {
// 					flp := flpc.GetLandmarkPoint(leftEye, rightEye, *imgParams, puploc.Perturbs, false)
// 					if flp.Row > 0 && flp.Col > 0 {
// 						dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
// 					}

// 					flp = flpc.GetLandmarkPoint(leftEye, rightEye, *imgParams, puploc.Perturbs, true)
// 					if flp.Row > 0 && flp.Col > 0 {
// 						dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
// 					}
// 				}
// 			}

// 			// Traverse all the mouth cascades and run the detector on each of them.
// 			for _, mouth := range mobiLM.mouthCascade {
// 				fmt.Println("-------- mount point-------")
// 				fmt.Println(len(mobiLM.flpcs[mouth]))
// 				for _, flpc := range mobiLM.flpcs[mouth] {
// 					flp := flpc.GetLandmarkPoint(leftEye, rightEye, *imgParams, puploc.Perturbs, false)
// 					if flp.Row > 0 && flp.Col > 0 {
// 						dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
// 					}
// 				}
// 			}
// 			flp := mobiLM.flpcs["lp84"][0].GetLandmarkPoint(leftEye, rightEye, *imgParams, puploc.Perturbs, true)
// 			if flp.Row > 0 && flp.Col > 0 {
// 				dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
// 				mobiLM.NoseCol = flp.Col
// 				mobiLM.NoseRow = flp.Row
// 			}
// 		}

// 		fmt.Println(len(dets[0]))
// 	}
// }

func (mobiLM *MobilePigoLandMark) GetFaceLandMark(bytesImage []uint8) {

	//// ----------------------------------------processing face Image
	var img, err = jpeg.Decode(bytes.NewReader(bytesImage))
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}
	fmt.Println("start")

	var pixels = pigo.RgbToGrayscale(img)
	fmt.Println(len(pixels))
	cols, rows := img.Bounds().Max.X, img.Bounds().Max.Y

	imgParams := &pigo.ImageParams{
		Pixels: pixels,
		Rows:   rows,
		Cols:   cols,
		Dim:    cols,
	}
	cParams := pigo.CascadeParams{
		MinSize:     60,
		MaxSize:     600,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: *imgParams,
	}
	//// ----------------------------------------processing face Detect

	filterResult := mobiLM.faceClassifier.RunCascade(cParams, 0.0)

	// Calculate the intersection over union (IoU) of two clusters.
	filterResult = mobiLM.faceClassifier.ClusterDetections(filterResult, 0.0)
	if len(filterResult) > 0 {
		results := []pigo.Detection{filterResult[0]}

		mobiLM.Scale = filterResult[0].Scale
		mobiLM.Row = filterResult[0].Row
		mobiLM.Col = filterResult[0].Col
		mobiLM.Rows = rows
		mobiLM.Cols = cols
		mobiLM.Q = int(filterResult[0].Q)

		dets := make([][]int, len(results))

		for i := 0; i < len(results); i++ {
			dets[i] = append(dets[i], results[i].Row, results[i].Col, results[i].Scale, int(results[i].Q), 0)
			chLeftEye := make(chan int, 3)
			chLeftRight := make(chan int, 3)
			chPuploc := make(chan pigo.Puploc, 1)
			chLeftEyePuploc := make(chan pigo.Puploc, 1)
			chRightEyePuploc := make(chan pigo.Puploc, 1)
			var wg sync.WaitGroup
			var puploc pigo.Puploc
			var leftEyePuploc pigo.Puploc
			var rightEyePuploc pigo.Puploc
			wg.Add(2)
			// left eye
			go func() {
				puploc := &pigo.Puploc{
					Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
					Col:      results[i].Col - int(0.185*float32(results[i].Scale)),
					Scale:    float32(results[i].Scale) * 0.4,
					Perturbs: 63,
				}
				leftEye := mobiLM.puplocClassifier.RunDetector(*puploc, *imgParams, 0.0, false)
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
				rightEye := mobiLM.puplocClassifier.RunDetector(*puploc, *imgParams, 0.0, false)
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
			puploc = <-chPuploc
			leftEyePuploc = <-chLeftEyePuploc
			rightEyePuploc = <-chRightEyePuploc

			// phase 2
			wg.Add(3)
			go func() {
				fmt.Println("-------- eye point-------")
				fmt.Println(len(mobiLM.eyeCascades))
				for _, eye := range mobiLM.eyeCascades {
					for _, flpc := range mobiLM.flpcs[eye] {
						flp := flpc.GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, false)
						if flp.Row > 0 && flp.Col > 0 {
							dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
						}
						flp = flpc.GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, true)
						if flp.Row > 0 && flp.Col > 0 {
							dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
						}
					}
				}
				wg.Done()
			}()
			go func() {
				fmt.Println("-------- eye point-------")
				fmt.Println(len(mobiLM.mouthCascade))
				for _, mouth := range mobiLM.mouthCascade {
					for _, flpc := range mobiLM.flpcs[mouth] {
						flp := flpc.GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, false)
						if flp.Row > 0 && flp.Col > 0 {
							dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
						}
					}
				}
				wg.Done()
			}()
			go func() {
				flp := mobiLM.flpcs["lp84"][0].GetLandmarkPoint(&leftEyePuploc, &rightEyePuploc, *imgParams, puploc.Perturbs, true)
				if flp.Row > 0 && flp.Col > 0 {
					dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
				}
				wg.Done()
			}()
			wg.Wait()

			// Traverse all the eye cascades and run the detector on each of them.

			// Traverse all the mouth cascades and run the detector on each of them.

		}

		fmt.Println(len(dets[0]))
	}
}
