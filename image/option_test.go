package image

import "testing"

func TestBimgOptions(t *testing.T) {
	imgOpts := imageOptions{
		Width:  480,
		Height: 620,
	}
	opts := bimgOptions(imgOpts)

	if opts.Width != imgOpts.Width || opts.Height != imgOpts.Height {
		t.Error("Invalid width/height")
	}
}
