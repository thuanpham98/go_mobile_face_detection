# go_mobile_face_detection

A Face Detection module write with Go language and binding for IOS and Android for mobile by Go Mobile

cmd binding code :
gomobile bind -androidapi=21 -target android -o build/android/faceDetection.aar
gomobile bind -target=ios/arm64 -iosversion=12 -o build/ios/FaceDetection.xcframework
#required Xcode
xcodebuild -create-xcframework -framework ./ios-arm64/FaceDetection.framework -framework ./ios-arm64_x86_64-simulator/FaceDetection.framework -output FaceDetection.xcframework
