# go_mobile_face_detection

A Face Detection module write with Go language and binding for IOS and Android for mobile by Go Mobile

cmd binding code :
gomobile bind -androidapi=21 -target android -o faceDetection.aar
gomobile bind -target=ios/arm64 -iosversion=12 -o build/ios/FaceDetection.xcframework
#required Xcode
