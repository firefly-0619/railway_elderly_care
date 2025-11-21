package utils

import (
	"fmt"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"io"
	"math"
	"strings"
)

// GetAudioDuration 从音频文件中获取时长
// file: 文件对象
// format: 文件格式/扩展名,如.mp3,.wav
// 返回微秒值
func GetAudioDuration(file interface {
	io.Reader
	io.Seeker
	io.Closer
}, format string) (float64, error) {
	// 重置文件指针到开头
	file.Seek(0, 0)
	format = strings.ToLower(format)
	defer file.Seek(0, 0)

	switch format {
	case ".mp3":
		streamer, format, err := mp3.Decode(file)
		if err != nil {
			return 0, err
		}
		defer streamer.Close()
		return float64(format.SampleRate.D(streamer.Len())), nil

	case ".wav":
		streamer, format, err := wav.Decode(file)
		if err != nil {
			return 0, err
		}
		defer streamer.Close()
		return float64(format.SampleRate.D(streamer.Len())), nil

	case ".flac":
		streamer, format, err := flac.Decode(file)
		if err != nil {
			return 0, err
		}
		defer streamer.Close()
		return float64(format.SampleRate.D(streamer.Len())), nil

	default:
		return 0, nil
	}
}

func FormatDurationStr(duration float64) string {
	round := int(math.Round(duration))
	minute := round / 60
	seconds := round % 60
	return fmt.Sprintf("%02d:%02d", minute, seconds)
}
