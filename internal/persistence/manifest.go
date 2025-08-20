package persistence

// Manifest JSON structure, keep track of the latest Snapshot and the AOF sequence
type Manifest struct {
	Version       int    `json:"version"`
	LastSnapshot  string `json:"lastSnapshot"`
	CurrentAOFKey string `json:"currentAOFKey"`
	AOFSeq        uint64 `json:"aofSequence"`
}
