package env

type DiffType string

const (
	DiffAdded     DiffType = "ADDED"
	DiffModified  DiffType = "MODIFIED"
	DiffRemoved   DiffType = "REMOVED"
	DiffUnchanged DiffType = "UNCHANGED"
)

type DiffItem struct {
	Key         string   `json:"key"`
	LocalValue  string   `json:"local_value,omitempty"`
	RemoteValue string   `json:"remote_value,omitempty"`
	Type        DiffType `json:"type"`
}

type DiffResult struct {
	Items      []DiffItem `json:"items"`
	HasChanges bool       `json:"has_changes"`
	Added      int        `json:"added"`
	Modified   int        `json:"modified"`
	Removed    int        `json:"removed"`
	Unchanged  int        `json:"unchanged"`
}

// CompareMaps compares two key-value maps (e.g. local vs remote envs)
func CompareMaps(local, remote map[string]string) *DiffResult {
	result := &DiffResult{Items: []DiffItem{}}

	visited := make(map[string]bool)

	// Check local items against remote
	for key, localVal := range local {
		visited[key] = true
		remoteVal, exists := remote[key]

		if !exists {
			result.Items = append(result.Items, DiffItem{
				Key:        key,
				LocalValue: localVal,
				Type:       DiffAdded,
			})
			result.Added++
			result.HasChanges = true
		} else if localVal != remoteVal {
			result.Items = append(result.Items, DiffItem{
				Key:         key,
				LocalValue:  localVal,
				RemoteValue: remoteVal,
				Type:        DiffModified,
			})
			result.Modified++
			result.HasChanges = true
		} else {
			result.Items = append(result.Items, DiffItem{
				Key:         key,
				LocalValue:  localVal,
				RemoteValue: remoteVal,
				Type:        DiffUnchanged,
			})
			result.Unchanged++
		}
	}

	// Check remote items missing locally
	for key, remoteVal := range remote {
		if !visited[key] {
			result.Items = append(result.Items, DiffItem{
				Key:         key,
				RemoteValue: remoteVal,
				Type:        DiffRemoved,
			})
			result.Removed++
			result.HasChanges = true
		}
	}

	return result
}

// CompareEnvs compares local EnvFile against remote map
func CompareEnvs(local *EnvFile, remote map[string]string) *DiffResult {
	localMap := local.Map()
	return CompareMaps(localMap, remote)
}
