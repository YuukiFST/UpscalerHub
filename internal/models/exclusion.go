package models

// ScanExclusion defines a rule for excluding games from scan results.
type ScanExclusion struct {
	Name        string `json:"Name"`
	PathSegment string `json:"PathSegment"`
}
