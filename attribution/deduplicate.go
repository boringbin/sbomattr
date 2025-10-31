package attribution

// Deduplicate removes duplicate attributions based on Purl, falling back to Name.
// The first occurrence of each unique attribution is kept.
func Deduplicate(attributions []Attribution) []Attribution {
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
		} else {
			logger.Debug("skipping duplicate attribution", "key", key) //nolint:sloglint // Package-level logger
		}
	}

	return result
}
