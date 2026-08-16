package env

type ItemType string

const (
	ItemEntry     ItemType = "ENTRY"
	ItemComment   ItemType = "COMMENT"
	ItemEmptyLine ItemType = "EMPTY"
)

type EnvItem struct {
	Type    ItemType `json:"type"`
	Key     string   `json:"key,omitempty"`
	Value   string   `json:"value,omitempty"`
	RawLine string   `json:"raw_line"`
	Comment string   `json:"comment,omitempty"`
}

type EnvFile struct {
	Items []*EnvItem `json:"items"`
}

func (f *EnvFile) Map() map[string]string {
	m := make(map[string]string)
	for _, item := range f.Items {
		if item.Type == ItemEntry {
			m[item.Key] = item.Value
		}
	}
	return m
}

func (f *EnvFile) Get(key string) (string, bool) {
	for _, item := range f.Items {
		if item.Type == ItemEntry && item.Key == key {
			return item.Value, true
		}
	}
	return "", false
}

func (f *EnvFile) Set(key, value string) {
	for _, item := range f.Items {
		if item.Type == ItemEntry && item.Key == key {
			item.Value = value
			return
		}
	}
	f.Items = append(f.Items, &EnvItem{
		Type:  ItemEntry,
		Key:   key,
		Value: value,
	})
}
