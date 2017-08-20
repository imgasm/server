package image

import "gopkg.in/h2non/bimg.v1"

type imageOptions struct {
	Width       int
	Height      int
	AreaWidth   int
	AreaHeight  int
	Quality     int
	Compression int
	Rotate      int
	Top         int
	Left        int
	Margin      int
	Factor      int
	DPI         int
	TextWidth   int
	Flip        bool
	Flop        bool
	Force       bool
	Embed       bool
	NoCrop      bool
	NoReplicate bool
	NoRotation  bool
	NoProfile   bool
	Opacity     float32
	Text        string
	Font        string
	Type        string
	Color       []uint8
	Extend      bimg.Extend
	Gravity     bimg.Gravity
	Colorspace  bimg.Interpretation
	Background  []uint8
}

func bimgOptions(imgOpts imageOptions) bimg.Options {
	opts := bimg.Options{
		Width:          imgOpts.Width,
		Height:         imgOpts.Height,
		Flip:           imgOpts.Flip,
		Flop:           imgOpts.Flop,
		Quality:        imgOpts.Quality,
		Compression:    imgOpts.Compression,
		NoAutoRotate:   imgOpts.NoRotation,
		NoProfile:      imgOpts.NoProfile,
		Force:          imgOpts.Force,
		Gravity:        imgOpts.Gravity,
		Embed:          imgOpts.Embed,
		Extend:         imgOpts.Extend,
		Interpretation: imgOpts.Colorspace,
		Type:           imageType(imgOpts.Type),
		Rotate:         bimg.Angle(imgOpts.Rotate),
	}

	if len(imgOpts.Background) != 0 {
		opts.Background = bimg.Color{imgOpts.Background[0], imgOpts.Background[1], imgOpts.Background[2]}
	}

	return opts
}
