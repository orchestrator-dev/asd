package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"asd/exec"
	"asd/render"
)

type VideoHandler struct {
	Runner exec.Runner
}

func (h *VideoHandler) CanHandle(mime, ext string) bool {
	ext = strings.ToLower(ext)
	switch ext {
	case ".mp4", ".mkv", ".avi", ".mov", ".webm", ".flv":
		return true
	}
	return strings.HasPrefix(mime, "video/")
}

func (h *VideoHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	out := render.NewWriter(w, opts.Flat)
	defer out.Close()

	if opts.Flat {
		_, err := io.Copy(out, r)
		return err
	}

	if h.Runner == nil {
		h.Runner = &exec.OSExecRunner{}
	}

	fmt.Fprintf(out, "File: %s\n", meta.Name)
	fmt.Fprintf(out, "Size: %d bytes\n", meta.Size)

	codec := "unknown"
	var width, height int
	fps := "unknown"
	duration := "unknown"
	bitrate := "unknown"

	if h.Runner.Available("ffprobe") && meta.Name != "stdin" && meta.Name != "" {
		output, err := h.Runner.Run("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", meta.Name)
		if err == nil {
			var probe struct {
				Format struct {
					Duration string `json:"duration"`
					BitRate  string `json:"bit_rate"`
				} `json:"format"`
				Streams []struct {
					CodecName  string `json:"codec_name"`
					CodecType  string `json:"codec_type"`
					Width      int    `json:"width"`
					Height     int    `json:"height"`
					RFrameRate string `json:"r_frame_rate"`
				} `json:"streams"`
			}
			if err := json.Unmarshal(output, &probe); err == nil {
				for _, s := range probe.Streams {
					if s.CodecType == "video" {
						codec = s.CodecName
						width = s.Width
						height = s.Height
						fps = s.RFrameRate
						break
					}
				}
				if probe.Format.Duration != "" {
					duration = probe.Format.Duration
				}
				if probe.Format.BitRate != "" {
					bitrate = probe.Format.BitRate
				}
			}
		}
	}

	fmt.Fprintf(out, "Codec: %s\n", codec)
	fmt.Fprintf(out, "Resolution: %dx%d\n", width, height)
	fmt.Fprintf(out, "FPS: %s\n", fps)
	fmt.Fprintf(out, "Duration: %s\n", duration)
	fmt.Fprintf(out, "Bitrate: %s\n", bitrate)
	fmt.Fprintf(out, "\nPlay hint: mpv %s\n", meta.Name)

	return nil
}
