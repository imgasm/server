package image

import (
	"gopkg.in/h2non/bimg.v1"
	"strings"
)

func extractImageTypeFromMime(mime string) string {
	mime = strings.Split(mime, ";")[0]
	parts := strings.Split(mime, "/")
	if len(parts) < 2 {
		return ""
	}
	name := strings.Split(parts[1], "+")[0]
	return strings.ToLower(name)
}

func isImageMimeTypeSupported(mime string) bool {
	ext := extractImageTypeFromMime(mime)
	if ext == "xml" {
		ext = "svg"
	}
	return bimg.IsTypeNameSupported(ext)
}

func imageType(name string) bimg.ImageType {
	ext := strings.ToLower(name)
	if ext == "jpeg" {
		return bimg.JPEG
	}
	if ext == "png" {
		return bimg.PNG
	}
	if ext == "webp" {
		return bimg.WEBP
	}
	if ext == "tiff" {
		return bimg.TIFF
	}
	if ext == "gif" {
		return bimg.GIF
	}
	if ext == "svg" {
		return bimg.SVG
	}
	if ext == "pdf" {
		return bimg.PDF
	}
	return bimg.UNKNOWN
}

func getImageMimeType(imgType bimg.ImageType) string {
	if imgType == bimg.PNG {
		return "image/png"
	}
	if imgType == bimg.WEBP {
		return "image/webp"
	}
	if imgType == bimg.TIFF {
		return "image/tiff"
	}
	if imgType == bimg.GIF {
		return "image/gif"
	}
	if imgType == bimg.SVG {
		return "image/svg+xml"
	}
	if imgType == bimg.PDF {
		return "application/pdf"
	}
	if imgType == bimg.TIFF {
		return "image/tiff"
	}
	return "image/jpeg"
}
