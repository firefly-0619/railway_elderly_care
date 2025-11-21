package utils

import (
	"path"
	"strings"
)

func ExtractFileSuffix(fileName string) string {

	return path.Ext(fileName)
}

func IsImageFile(fileName string) bool {
	ext := path.Ext(fileName)
	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".png") {
		return true
	}
	return false
}

func IsMusicFile(fileName string) bool {
	ext := path.Ext(fileName)
	if strings.EqualFold(ext, ".mp3") || strings.EqualFold(ext, ".wav") || strings.EqualFold(ext, ".flac") {
		return true
	}
	return false
}

func IsLRCFile(fileName string) bool {
	ext := path.Ext(fileName)
	if strings.EqualFold(ext, ".lrc") {
		return true
	}
	return false
}
