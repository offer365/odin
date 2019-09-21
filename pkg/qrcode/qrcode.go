package qrcode

import (
	"bytes"
	"errors"
	qr "github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
)

const (
	MaxWidth int = 360 // 二维码图片的最大宽度
	MaxLang  int = 800 // 每个二维码图片的最大字节数
)

type Code struct {
	data []byte
}

//实现 io.Writer 接口
func (c *Code) Write(p []byte) (n int, err error) {
	c.data = p
	return len(c.data), nil
}

func NewQrEncode(content []byte, buf io.Writer) error {
	var FirstImg image.Image
	if len(content) == 0 {
		return errors.New("No data.")
	}
	bnr := bytes.NewReader(content)
	imgs := make([]image.Image, 0)
	height := 0
	for {
		data := make([]byte, MaxLang)
		_, err := io.ReadFull(bnr, data)
		if err == io.EOF {
			break
		}
		img := QrImg(string(data), MaxWidth)
		if img != nil {
			imgs = append(imgs, img)
		}
	}
	for _, i := range imgs {
		height += i.Bounds().Max.Y
	}
	NewImg := image.NewNRGBA(image.Rect(0, 0, MaxWidth, height)) //创建一个新RGBA图像
	if len(imgs) > 0 {
		FirstImg = imgs[0]
		draw.Draw(NewImg, NewImg.Bounds(), FirstImg, FirstImg.Bounds().Min, draw.Over) //画上第一张缩放后的图片
		if len(imgs) > 1 {
			for i, img := range imgs[1:] {
				draw.Draw(NewImg, NewImg.Bounds(), img, img.Bounds().Min.Sub(image.Pt(0, FirstImg.Bounds().Max.Y*(i+1))), draw.Over) //画上第二张缩放后的图片（这里需要注意Y值的起始位置）
			}
		}

		//jpg, _ := os.Create("qr-code.jpg")  //写入到文件
		//defer jpg.Close()
		//err:=jpeg.Encode(jpg, NewImg, &jpeg.Options{100})

		return jpeg.Encode(buf, NewImg, &jpeg.Options{Quality: jpeg.DefaultQuality + 25}) // &jpeg.Options{100} 图片质量最好
	}
	return errors.New("No image.")

}

func QrImg(content string, size int) (img image.Image) {
	data, err := qr.Encode(content, qr.Medium, size)
	if err != nil {
		return nil
	}
	buf := bytes.NewBuffer(data)
	if img, _, err = image.Decode(buf); err != nil {
		return nil
	}
	return
}

func NewQrCode(content string, size int, path string) error {
	return qr.WriteFile(content, qr.Medium, size, path)
}
