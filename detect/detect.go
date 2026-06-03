package detect

// Detector defines an interface for detecting file types.
type Detector interface {
	Detect(header []byte, filename string) (mime, confidence string)
}

// Pipeline runs a sequence of detectors to identify a file.
func Pipeline(header []byte, filename string) string {
	detectors := []Detector{
		&MagicDetector{},
		&ExtDetector{},
		&ShebangDetector{},
	}

	bestMime := "application/octet-stream"

	for _, d := range detectors {
		mime, conf := d.Detect(header, filename)
		if conf == "high" && mime != "" {
			return mime
		}
		if conf == "medium" && mime != "" {
			bestMime = mime
		}
	}

	return bestMime
}
