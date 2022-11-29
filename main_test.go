package imwrapper

import (
	"errors"
	"os"
	"strings"
	"testing"
)

type testCase struct {
	inputImage     string
	expectedWidth  int
	expectedHeight int
}
type convertTest struct {
	inputImage   string
	outputImage  string
	maxDimension int
	qualityParam int
	autoRotate   bool
	error        bool
	title        string
}

func TestCmd_HappyFlows(t *testing.T) {
	err, im := New()
	if err != nil {
		t.Fatalf("could not instantiate imagemagick: %v", err)
	}

	testCases := []testCase{
		{
			inputImage:     "testdata/input.jpg",
			expectedWidth:  1024,
			expectedHeight: 768,
		},
	}

	for _, test := range testCases {
		width, height, err := im.getDimensions(test.inputImage)
		if err != nil {
			t.Fatalf("Test failed, could not get image dimensions of file: %v", err)
		}
		if width != test.expectedWidth {
			t.Fatalf("Image does not have expected width: %v, but instead [%v]", test.expectedWidth, width)
		}
		if height != test.expectedHeight {
			t.Fatalf("Image does not have expected height: %v, but instead [%v]", test.expectedHeight, height)
		}
	}
}

func TestCmd_Convert(t *testing.T) {
	err, im := New()
	if err != nil {
		t.Fatalf("could not instantiate imagemagick: %v", err)
	}

	testCases := []convertTest{
		{
			inputImage:   "testdata/input.jpg",
			outputImage:  "testdata/output.jpg",
			maxDimension: 300,
			qualityParam: 5,
			autoRotate:   false,
			error:        false,
			title:        "simple happy flow",
		},
		{
			inputImage:   "testdata/input.jpg",
			outputImage:  "testdata/output.jpg",
			maxDimension: 32,
			qualityParam: 55,
			autoRotate:   true,
			error:        false,
			title:        "pass auto-orient parameter",
		},
		{
			inputImage:  "testdata/input.jpg",
			outputImage: "testdata/input.jpg",
			error:       true,
			title:       "source should not equal destination",
		},
	}

	for _, test := range testCases {
		err := im.Convert(test.inputImage, test.outputImage, test.qualityParam, test.maxDimension, test.autoRotate)
		if test.error && err == nil || !test.error && err != nil {
			t.Fatalf("Test failed: %v", test.title)
		} else {
			cleanup(test.outputImage)
		}
	}
}

func cleanup(image string) {
	if strings.HasSuffix(image, "input.jpg") {
		return
	}
	if _, err := os.Stat(image); !errors.Is(err, os.ErrNotExist) {
		os.Remove(image)
	}
}

func TestCheckFiles_happy(t *testing.T) {
	err := checkFiles("./main.go", "./main")
	if err != nil {
		t.Fatalf("System checks did not pass due to error: %v", err)
	}
}

func TestCheckFiles_no_source(t *testing.T) {
	err := checkFiles("./some-non-existing-file.txt", "./main")
	if err == nil {
		t.Fatalf("System checks did not throw an error for missing source...")
	}
}

func TestCheckFiles_no_destination(t *testing.T) {
	err := checkFiles("./main.go", "./main.go")
	if err == nil {
		t.Fatalf("System checks did not throw an error for existing destination...")
	}
}
