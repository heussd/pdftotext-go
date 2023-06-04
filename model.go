package pdftotext

type PdfPage struct {
	Content string
	// PDF page number (1-based)
	Number int
}

type PopplerTsvRow struct {
	Level    int     `col:"0"`
	PageNum  int     `col:"1"`
	ParNum   int     `col:"2"`
	BlockNum int     `col:"3"`
	LineNum  int     `col:"4"`
	WordNum  int     `col:"5"`
	Left     float64 `col:"6"`
	Top      float64 `col:"7"`
	Width    float64 `col:"8"`
	Height   float64 `col:"9"`
	Conf     int     `col:"10"`
	Text     string  `col:"11"`
}
