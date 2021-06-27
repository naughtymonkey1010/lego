package image

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestWhatImpl(t *testing.T) {

	//file,_:=os.Open("/Users/admin/Downloads/0SVgfloXTP.webp")
	//file,_:=os.Open("/Users/admin/Downloads/0SGo2bnbwT.png")
	//file, _ := os.Open("/Users/admin/Downloads/0Me1EWJKUy.jpeg")
	//file, _ := os.Open("/Users/admin/Downloads/0TLSqPSaGU.gif")
	//file, _ := os.Open("/Users/admin/Downloads/图片上传流程.jpg")
	//file, _ := os.Open("/Users/admin/Downloads/0oz4MyLoAL2.jpeg")

	//file, _ := os.Open("/Users/admin/Downloads/0p8raLrd8F5.jpeg")
	//file, _ := os.Open("/Users/admin/Downloads/WechatIMG19.jpeg")

	//file, _ := os.Open("/Users/admin/Downloads/WechatIMG9903.jpeg")
	//file, _ := os.Open("/Users/admin/Downloads/WechatIMG9904.jpeg")
	//file, _ := os.Open("/Users/admin/Downloads/WechatIMG9905.jpeg")
	//file, _ := os.Open("/Users/admin/Downloads/WechatIMG9906.jpeg")
	//file, _ := os.Open("/Users/admin/Downloads/WechatIMG9910.jpeg")
	//file, _ := os.Open("/Users/admin/Downloads/WechatIMG9913.jpeg")

	file, _ := os.Open("/Users/admin/Documents/IMG_3019.JPG") //6
	//file, _ := os.Open("/Users/admin/Documents/IMG_3020.JPG")//3
	//file, _ := os.Open("/Users/admin/Documents/IMG_3021.JPG")  //8
	//file, _ := os.Open("/Users/admin/Documents/IMG_3022.JPG")

	//file, _ := os.Open("/Users/admin/Downloads/0pjIJmHKHfC.jpeg")

	// rename file
	//file, _ := os.Open("/Users/admin/Downloads/0SGo2bnbwT.webp")
	//format, err := WhatImpl(file)
	//fmt.Printf("format:%s    err:%s", format, err)
	//file, _ := os.Open("/Users/admin/Downloads/test-webp-animated.webp")
	//file, _ := os.Open("/Users/admin/Downloads/2.webp")
	//file, _ := os.Open("/Users/admin/Downloads/git-animated.gif")

	//image.RegisterFormat("webp", "RIFF????WEBPVP8", webp.Decode, webp.DecodeConfig)
	//image.RegisterFormat("gif", "GIF8?a", gif.Decode, gif.DecodeConfig)
	//image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)
	//image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)

	////fmt.Println("===================")
	//config,err:=webp.DecodeConfig(file)
	//if err!=nil {
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(config.Height)
	//	fmt.Println(config.Width)
	//}
	////
	//fmt.Println("===================")
	//

	//img, _, err := image.Decode(file)
	//
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(img.Bounds().Dx())
	//	fmt.Println(img.Bounds().Dy())
	//	//fmt.Println(name)
	//	fmt.Println(img.Bounds().Size())
	//}
	//
	//fmt.Println("===================")
	byteArr, _ := ioutil.ReadAll(file)
	fmt.Println(len(byteArr))
	//

	//
	//startTime2:=time.Now().Unix()
	//RotateImage(byteArr,6)
	//endTime2:=time.Now().Unix()
	//fmt.Println(endTime2-startTime2)

	//
	//startTime3:=time.Now().Unix()
	//RotateImage2(byteArr,270)
	//endTime3:=time.Now().Unix()
	//fmt.Println(endTime3-startTime3)

	info, _, _ := GetImageInfo(byteArr)
	fmt.Println(info)

}
