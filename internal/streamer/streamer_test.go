package streamer

import (
	"testing"
)

func TestEncodeToHLS(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}
	
	// Create a video that converts mp4 to web ready format.
	video1 := NewVideo(
		1, 
		"/mnt/d/workspace/src/domain/social/scripts/input/example.mp4", 
		"/mnt/d/workspace/src/domain/social/scripts/output",
	)

	err := video1.EncodeToHLS()
	
	if err != nil {
		t.Fatal(err)
	}

	t.Log("视频转换成功")
}


// todo "./scripts/input/bad.txt", "./output", 
