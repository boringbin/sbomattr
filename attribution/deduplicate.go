package attribution

import "log/slog"

// Deduplicate removes duplicate attributions based on Purl, falling back to Name.
// The first occurrence of each unique attribution is kept.
// The logger parameter is optional; pass nil to disable logging.
func Deduplicate(attributions []Attribution, logger *slog.Logger) []Attribution {
	seen := make(map[string]bool)
	result := make([]Attribution, 0, len(attributions))

	for _, a := range attributions {
		// Use Purl as primary key, fall back to Name if Purl is empty
		key := a.Purl
		if key == "" {
			key = a.Name
		}

		if !seen[key] {
			seen[key] = true
			result = append(result, a)
		} else if logger != nil {
			logger.Debug("skipping duplicate attribution", "key", key)
		}
	}

	return result
}
