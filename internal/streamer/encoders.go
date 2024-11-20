package streamer

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/xfrr/goffmpeg/transcoder"
)

// Encoder 视频编码的接口
type Encoder interface {
	EncodeToMP4(v *Video, baseFileName string) error
	EncodeToHLS(v *Video, baseFileName string) error
}

type VideoEncoder struct{}

// Takes a video object and a base file name and encodes to mp4
// 将视频编码为 MP4 格式
func (ve *VideoEncoder) EncodeToMP4(v *Video, baseFileName string) error {
	// Create transcoder
	trans := new(transcoder.Transcoder)

	// 构建输出路径
	outputPath := fmt.Sprintf("%s/%s.mp4", v.OutputDir, baseFileName)

	// Initialize the transcoder
	err := trans.Initialize(v.InputFile, outputPath)
	if err != nil {
		fmt.Println("Error encoding....")
		return err
	}

	// Set Codec.
	trans.MediaFile().SetVideoCodec("libx264")

	// Start the transcoding process
	done := trans.Run(false)

	err = <-done
	if err != nil {
		return err
	}

	return nil
}

func (vd *VideoEncoder) EncodeToHLS(v *Video, baseFileName string) error {
	// Create a channel to get results
	result := make(chan error)

	// Spawn a goroutine to do the encode
	go func(result chan error) {
		ffmpegCmd := exec.Command(
			"ffmpeg",
			"-i", v.InputFile,
			"-map", "0:v:0",
			"-map", "0:a:0",
			"-map", "0:v:0",
			"-map", "0:a:0",
			"-map", "0:v:0",
			"-map", "0:a:0",
			"-c:v", "libx264",
			"-crf", "22",
			"-c:a", "aac",
			"-ar", "48000",
			"-filter:v:0", "scale=-2:1080",
			"-maxrate:v:0", v.Options.MaxRate1080p,
			"-b:a:0", "128k",
			"-filter:v:1", "scale=-2:720",
			"-maxrate:v:1", v.Options.MaxRate720p,
			"-b:a:1", "128k",
			"-filter:v:2", "scale=-2:480",
			"-maxrate:v:2", v.Options.MaxRate480p,
			"-b:a:2", "64k",
			"-var_stream_map", "v:0,a:0,name:1080p v:1,a:1,name:720p v:2,a:2,name:480p",
			"-preset", "slow",
			"-hls_list_size", "0",
			"-threads", "0",
			"-f", "hls",
			"-hls_playlist_type", "event",
			"-hls_time", strconv.Itoa(v.Options.SegmentDuration),
			"-hls_flags", "independent_segments",
			"-hls_segment_type", "mpegts",
			"-hls_playlist_type", "vod",
			"-master_pl_name", fmt.Sprintf("%s.m3u8", baseFileName),
			"-profile:v", "baseline",
			"-level", "3.0",
			"-progress", "-",
			"-nostats",
			fmt.Sprintf("%s/%s-%%v.m3u8", v.OutputDir, baseFileName),
		)

		// ffmpeg -i input.mp4 -c:v libx264 -c:a aac -strict -2 -f hls -hls_time 10 -hls_list_size 0 output.m3u8
		// ffmpeg -i example.mp4 -c:v libx264 -c:a aac -strict -2 -f hls -hls_time 10 -hls_list_size 0 output.m3u8
		
		_, err := ffmpegCmd.CombinedOutput()
		result <- err
	}(result)

	// Listen to the result channel
	err := <-result
	if err != nil {
		return err
	}

	// Return the results (success or not)
	return nil
}
