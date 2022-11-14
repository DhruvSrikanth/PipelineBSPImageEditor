package task

type Task struct {
	InPath  string   `json:"inPath"`
	OutPath string   `json:"outPath"`
	Effects []string `json:"effects"`
}
