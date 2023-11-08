package main

type Media struct {
	Url      string
	FilePath string
}

type Download struct {
	Url       string
	FilePath  string
	Error     error
	IsSuccess bool
	Hash      string
}

// region - Sort

type ByFilePath []Download

func (a ByFilePath) Len() int           { return len(a) }
func (a ByFilePath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFilePath) Less(i, j int) bool { return a[i].FilePath < a[j].FilePath }

// endregion
