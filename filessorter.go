package main

import (
	"os"
)

type FilesSorter struct {
	files []os.FileInfo
}

func (s *FilesSorter) Len() int {
	return len(s.files)
}

func (s *FilesSorter) Swap(i, j int) {
	s.files[i], s.files[j] = s.files[j], s.files[i]
}

func (s *FilesSorter) Less(i, j int) bool {
	return s.files[i].Name() < s.files[j].Name()
}
