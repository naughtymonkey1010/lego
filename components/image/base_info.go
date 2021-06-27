package image

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"golang.org/x/image/webp"
)

type ImgInfo struct {
	Format   string
	Allow    bool
	Size     [2]int
	Animated bool
}

func init() {
	//注册图片格式解析器
	image.RegisterFormat("webp", "RIFF????WEBPVP8", webp.Decode, webp.DecodeConfig)
	image.RegisterFormat("gif", "GIF8?a", gif.Decode, gif.DecodeConfig)
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
}

func GetImageInfo(byteArr []byte) (*ImgInfo, []byte, error) {

	//
	imgInfo := &ImgInfo{}
	img, format, err := image.Decode(bytes.NewReader(byteArr))
	//如果图片是jpeg格式  需要考虑图片的旋转问题 获取宽高等信息异常  需要进行旋转
	if strings.ToLower(format) == "jpeg" {
		orientVal, _ := ReadOrientation(byteArr)
		if orientVal == 3 || orientVal == 6 || orientVal == 8 {
			byteArr, _ = RotateImage(byteArr, orientVal)
			//旋转后图片要重新进行解析
			img, format, err = image.Decode(bytes.NewReader(byteArr))
		}
	}
	if err != nil && strings.Contains(err.Error(), "VP8X") {
		imgInfo.Format = WEBP
		imgInfo.Allow = true
		imgInfo.Animated = true
		imgInfo.Size = [2]int{-1, -1}
		return imgInfo, byteArr, nil
	} else if err != nil && strings.Contains(err.Error(), "webp") {
		imgInfo.Format = WEBP
		imgInfo.Allow = true
		imgInfo.Animated = false
		imgInfo.Size = [2]int{-1, -1}
		return imgInfo, byteArr, nil
	} else if err != nil {
		return nil, byteArr, err
	}
	imgInfo.Format = format
	imgInfo.Allow = true
	imgInfo.Animated = false
	imgInfo.Size = [2]int{img.Bounds().Dx(), img.Bounds().Dy()}
	if format == GIF {
		gifs, err := gif.DecodeAll(bytes.NewReader(byteArr))
		if err != nil {
			return nil, byteArr, err
		}
		imgInfo.Animated = len(gifs.Image) > 1
	}
	return imgInfo, byteArr, nil
}

func ReadOrientation(byteArr []byte) (int, error) {
	x, err := exif.Decode(bytes.NewReader(byteArr))
	if err != nil {
		return 0, err
	}
	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		return 0, err
	}
	orientVal, err := orientation.Int(0)
	if err != nil {
		return 0, err
	}
	return orientVal, nil
}

func RotateImage(src []byte, orientVal int) ([]byte, error) {
	var img, _, err = image.Decode(bytes.NewReader(src))
	if err != nil {
		return src, err
	}
	var dst *image.NRGBA
	switch orientVal {
	case 3:
		dst = imaging.Rotate180(img)
	case 6:
		dst = imaging.Rotate270(img)
	case 8:
		dst = imaging.Rotate90(img)
	}
	//file, err := os.Create("/Users/admin/Documents/test1.jpg")
	//if err != nil {
	//	return src, err
	//}
	//defer file.Close()
	//err = jpeg.Encode(file, dst, &jpeg.Options{50})
	//return nil,nil

	//返回新的图片的二进制
	var writer bytes.Buffer
	err = jpeg.Encode(&writer, dst, &jpeg.Options{50})
	if err != nil {
		return src, err
	}
	return writer.Bytes(), nil
}
