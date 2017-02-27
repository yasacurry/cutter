package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
)

type byCustom struct {
	f    []os.FileInfo
	less func(i, j os.FileInfo) bool
}

func (c byCustom) Len() int {
	return len(c.f)
}

func (c byCustom) Less(i, j int) bool {
	return c.less(c.f[i], c.f[j])
}

func (c byCustom) Swap(i, j int) {
	c.f[i], c.f[j] = c.f[j], c.f[i]
}

const srcDir = "srcdir"
const dstDir = "dstdir"

func main() {
	fileinfo, _ := ioutil.ReadDir(srcDir)

	sort.Sort(byCustom{fileinfo, func(i, j os.FileInfo) bool {
		if i.ModTime().Unix() != j.ModTime().Unix() {
			return i.ModTime().Unix() < j.ModTime().Unix()
		}
		if i.Name() != j.Name() {
			return i.Name() < j.Name()
		}
		return false
	}})

	count := 0
	for _, f := range fileinfo {
		if isPng(f) == false {
			continue
		}
		count++
		srcPath := filepath.Join(".", srcDir, f.Name())
		dstName := fmt.Sprintf("%03d", count) + ".png"
		dstPath := filepath.Join(".", dstDir, dstName)
		cut(srcPath, dstPath, image.Point{1250, 420}, image.Point{2330, 1500})
	}
}

func cut(srcPath string, dstPath string, topLeft image.Point, bottomRight image.Point) {
	src, _ := os.Open(srcPath)
	defer src.Close()
	srcPng, _ := png.Decode(src)

	minP := image.Point{0, 0}
	maxP := image.Point{bottomRight.X - topLeft.X, bottomRight.Y - topLeft.Y}
	out := image.NewRGBA(image.Rectangle{minP, maxP})

	for y := out.Rect.Min.Y; y < out.Rect.Max.Y; y++ {
		for x := out.Rect.Min.X; x < out.Rect.Max.X; x++ {
			c := srcPng.At(x+topLeft.X, y+topLeft.Y)
			out.Set(x, y, c)
		}
	}

	dst, _ := os.Create(dstPath)
	defer dst.Close()
	png.Encode(dst, out)
}

// func cut(srcPath string, dstPath string, topLeft image.Point, bottomRight image.Point) {
// 	fmt.Println(srcPath, dstPath)
// }

func isPng(f os.FileInfo) bool {
	if ".png" == path.Ext(f.Name()) {
		return true
	}
	return false
}
