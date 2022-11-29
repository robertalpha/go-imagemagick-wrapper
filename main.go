package imwrapper

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type ImageMagick struct {
	convertPath  string
	identifyPath string
}

func New() (error, *ImageMagick) {
	identify, err := exec.LookPath("identify")
	if err != nil {
		return getError("could not find 'identify' command, please make sure imagemagick installed", ""), nil
	}
	convert, err := exec.LookPath("convert")
	if err != nil {
		return getError("could not find 'convert' command, please make sure imagemagick installed", ""), nil
	}
	return nil, &ImageMagick{
		identifyPath: identify,
		convertPath:  convert,
	}
}

func (i *ImageMagick) Convert(sourcePath string, destinationPath string, quality int, maxDimension int, autorotate bool) error {
	if err := checkFiles(sourcePath, destinationPath); err != nil {
		return err
	}

	if quality < 1 || quality > 100 {
		return getError("quality should be between 1 and 100", quality)
	}

	widthParam, heightParam, err := i.getDimensions(sourcePath)
	if err != nil {
		return err
	}
	width := getMin(widthParam, maxDimension)
	height := getMin(heightParam, maxDimension)

	params := formatImagemagickParams(width, height, autorotate, sourcePath, destinationPath)
	cmd := exec.Command(i.convertPath, strings.Split(params, " ")...)

	if err := cmd.Run(); err != nil {
		return getError("Error: %v", err)
	}
	return nil
}

func (i *ImageMagick) getDimensions(sourcePath string) (int, int, error) {
	cmd := exec.Command(i.identifyPath, "-ping", "-format", "%w %h", sourcePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, getError("Error: %v", err)
	}
	if len(output) == 0 {
		return 0, 0, getError("could not determine image dimensions", "")
	}
	dims := strings.Split(string(output), " ")
	if len(dims) != 2 {
		return 0, 0, getError("could not determine image dimensions", "")
	}
	width, err := strconv.Atoi(dims[0])
	if err != nil {
		return 0, 0, getError("could not determine image dimensions", err)
	}
	height, err := strconv.Atoi(dims[1])
	if err != nil {
		return 0, 0, getError("could not determine image dimensions", err)
	}

	return width, height, nil
}

func checkFiles(sourcePath string, destinationPath string) error {
	if sourcePath == destinationPath {
		return getError("destination should not equal source", destinationPath)
	}
	if _, err := os.Stat(sourcePath); errors.Is(err, os.ErrNotExist) {
		return getError("source file does not exist", destinationPath)
	}
	dstDir := path.Dir(destinationPath)
	if stat, err := os.Stat(dstDir); err == nil && !stat.IsDir() {
		return getError("destination directory does not exist", dstDir)
	}
	if _, err := os.Stat(destinationPath); err == nil {
		return getError("destination file already exists", destinationPath)
	}

	return nil
}

func getError(decription string, param interface{}) error {
	return fmt.Errorf("[go-imagemagick-wrapper] Error: %s [%v]", decription, param)
}

func formatImagemagickParams(witdh int, height int, autoRotateEnabled bool, input string, output string) string {
	autoRotate := ""
	if autoRotateEnabled {
		autoRotate = " -auto-orient"
	}
	return fmt.Sprintf("-interlace Plane%s -sampling-factor 4:2:0 -resize %dx%d -strip -quality 60 %s %s", autoRotate, witdh, height, input, output)
}

func getMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}
