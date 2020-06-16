package cmd

type Filter struct {
	Language string `json:"language"`
	Regex    string `json:"regex"`
}

type Target struct {
	Outfile string   `json:"outfile"`
	Sources []string `json:"sources"`
	Filters *Filter  `json:"filters"`
	Order   []string `json:"order"`
}
