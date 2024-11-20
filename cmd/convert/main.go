package main

import (
	"fmt"

	"github.com/sjxiang/social/internal/streamer"
)

func main() {

	// ffmpeg -i input_video.mp4 -ss 0:0:5 -vframes 1 -q:v 2 output_cover.jpg

	// ffmpeg -i example.mp4 -ss 0:0:5 -vframes 1 -q:v 2 output_cover.jpg
	// Define number of workers and jobs
	const numJobs = 4
	const numWorkers = 4

	// Create 2 channels for work and results (1. notifications and 2. We send work to)
	notifyChan := make(chan streamer.ProcessingMessage, numJobs)
	defer close(notifyChan)

	videoQueue := make(chan streamer.VideoProcessingJob, numJobs)
	defer close(videoQueue)

	// Get a worker pool
	wp := streamer.New(videoQueue, numWorkers)

	// Start the worker pool
	wp.Run()
	fmt.Println("Worker Pool Started. Press enter to continue")
	_, _ = fmt.Scanln()

	// Create 1 video to send to the worker pool
	ops := &streamer.VideoOptions{
		RenameOutput:    true,
		SegmentDuration: 10,
		MaxRate1080p:    "1200k",
		MaxRate720p:     "600k",
		MaxRate480p:     "400k",
	}

	// Create a video that converts mp4 to web ready format.
	video1 := wp.NewVideo(1, "./scripts/input/example.mp4", "./output", "mp4", notifyChan, nil)

	// Create second video that should fail.
	video2 := wp.NewVideo(2, "./scripts/input/bad.txt", "./output", "mp4", notifyChan, nil)

	// Create a third video that should convert mp4 to hls
	video3 := wp.NewVideo(3, "./scripts/input/example2.mp4", "./output", "hls", notifyChan, ops)

	// Create a fourth video that should convert mp4 to hls
	video4 := wp.NewVideo(4, "./scripts/input/example2.mp4", "./output", "mp4", notifyChan, nil)

	// Send the videos to the worker pool
	videoQueue <- streamer.VideoProcessingJob{
		Video: video1,
	}
	videoQueue <- streamer.VideoProcessingJob{
		Video: video2,
	}
	videoQueue <- streamer.VideoProcessingJob{
		Video: video3,
	}
	videoQueue <- streamer.VideoProcessingJob{
		Video: video4,
	}

	// Print out the results
	for i := 1; i <= numJobs; i++ {
		msg := <-notifyChan
		fmt.Println("i:", i, "/", "message:", msg)
	}

	fmt.Println("Done")
}
