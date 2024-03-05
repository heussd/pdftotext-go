# pdftotext-go

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/heussd/pdftotext-go/badge)](https://securityscorecards.dev/viewer/?uri=github.com/heussd/pdftotext-go)

Extract texts with their corresponding page numbers from PDF files.
Wraps the command line tool `pdftotext` ([poppler-utils](https://poppler.freedesktop.org/)).

## Usage

1. [poppler-utils](https://poppler.freedesktop.org/) (version >=22.05.0) must be installed and available in the path.
1. `go get "github.com/heussd/pdftotext-go"`
1. See [tests for code examples](pdftotext_test.go).



## Why poppler version >=22.05.0

[Version 22.05.0 of poppler](https://poppler.freedesktop.org/releases.html) introduced a new parameter `-tsv`, which extracts PDF content with meta data as TSV. This functionality is essential for the operation of this library.