package utils

import (
	"fmt"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func Test() {
	fmt.Println("init")
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	fmt.Println("read")
	err := mw.ReadImage("leezi/wut.png")
	if err != nil {
		panic(err)
	}
}
