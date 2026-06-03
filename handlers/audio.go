package handlers

import (
        "asd/exec"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/dhowden/tag"

	"asd/render"
)

type AudioHandler struct {
	Runner exec.Runner
}

func (h *AudioHandler) CanHandle(mime, ext string) bool {
	ext = strings.ToLower(ext)
	switch ext {
	case ".mp3", ".flac", ".ogg", ".wav", ".aac", ".m4a", ".opus":
		return true
	}
	return strings.HasPrefix(mime, "audio/")
}

func (h *AudioHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	out := render.NewWriter(w, opts.Flat)
	defer out.Close()

	if opts.Flat {
		_, err := io.Copy(out, r)
		return err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	m, err := tag.ReadFrom(bytes.NewReader(data))
	if err != nil {
		fmt.Fprintf(out, "File: %s\n", meta.Name)
		fmt.Fprintf(out, "Size: %d bytes\n", meta.Size)
		fmt.Fprintf(out, "Error reading metadata: %v\n", err)
		return nil
	}

	fmt.Fprintf(out, "File: %s\n", meta.Name)
	fmt.Fprintf(out, "Title: %s\n", m.Title())
	fmt.Fprintf(out, "Artist: %s\n", m.Artist())
	fmt.Fprintf(out, "Album: %s\n", m.Album())
	fmt.Fprintf(out, "Year: %d\n", m.Year())
	track, totalTracks := m.Track()
	if totalTracks > 0 {
		fmt.Fprintf(out, "Track: %d/%d\n", track, totalTracks)
	} else {
		fmt.Fprintf(out, "Track: %d\n", track)
	}

	duration := "unknown"
	bitrate := "unknown"

	fmt.Fprintf(out, "Duration: %s\n", duration)
	fmt.Fprintf(out, "Bitrate: %s\n", bitrate)

	return nil
}
