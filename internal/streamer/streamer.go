package streamer

import (
	"fmt"
	"os/exec"
	"strings"
	"path"
	"path/filepath"
)


type Encoder interface {
	EncodeToHLS(v *Video, baseFileName string) error
}


// GenerateCoverFromVideo


type Video struct {
	ID           int
	InputFile    string // the video we want to encode
	OutputDir    string // where we want the encoded video to show up
	BaseFileName string
}


func NewVideo(id int, inputFile, outputDir string) *Video {
	
	b := path.Base(inputFile)  // ex. '/input/cat.mp4' becomes 'cat.mp4'  
	
	return &Video{
		ID:           id,
		InputFile:    inputFile,
		OutputDir:    outputDir,
		BaseFileName: strings.TrimSuffix(b, filepath.Ext(b)),  // ex. cat.mp4 becomes cat, 移除后缀
	}
}

func (v *Video) EncodeToHLS() error {

	// Create a channel to get results
	result := make(chan error)

	// 例, ffmpeg -i input.mp4 -c:v libx264 -c:a aac -strict -2 -f hls -hls_time 10 -hls_list_size 0 output.m3u8
	
	// Spawn a goroutine to do the encode
	go func(result chan error) {
		ffmpegCmd := exec.Command(
			"ffmpeg",
			"-i", v.InputFile,
			"-c:v", "libx264",
			"-c:a", "aac",
			"-strict", "-2",
			"-f", "hls",
			"-hls_list_size", "0",
			"-hls_time", "10",
			fmt.Sprintf("%s/%s.m3u8", v.OutputDir, v.BaseFileName),  
		)

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
