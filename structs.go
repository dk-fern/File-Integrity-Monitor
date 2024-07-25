package main

type BaselineFileList struct {
	BaselinePath string `json:"baselinePath"`
	Files        []File `json:"files"`
}

type File struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

type Difference struct {
	HashDifferences []string
	AddedFiles      []string
	RemovedFiles    []string
}

type CompareFileList struct {
	BaselinePath string     `json:"baselinePath"`
	Differences  Difference `json:"differences"`
}
