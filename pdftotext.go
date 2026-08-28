package pdftotext

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

// https://poppler.freedesktop.org/releases.html
// poppler introduced the "tsv" parameter with this version:
const popplerVersionConstraint = ">= 22.05.0"

// Extract PDF text content in simplified format, it uses context.Background under the hood. Use ExtractContext to specify a context
func Extract(pdfBytes []byte) ([]PdfPage, error) {
	return ExtractContext(context.Background(), pdfBytes)
}

// ExtractContext behaves the same as Extract but lets the caller pass a Context
func ExtractContext(ctx context.Context, pdfBytes []byte) ([]PdfPage, error) {
	var pdfPages []PdfPage

	tsv, err := ExtractInPopplerTsvContext(ctx, pdfBytes)

	if err != nil {
		return nil, err
	}

	prevPage := 1
	prevContent := ""

	for i, row := range tsv {
		if row.Conf != -1 { // Seems to indicate control sequences
			prevContent += row.Text + " "
		}

		var (
			pageChanged   = prevPage != row.PageNum
			lastIteration = i == len(tsv)-1
		)

		if pageChanged || lastIteration {
			pdfPages = append(pdfPages, PdfPage{
				Content: prevContent,
				Number:  prevPage,
			})

			prevPage = row.PageNum
			prevContent = ""
		}
	}

	return pdfPages, nil
}

// ExtractOrError Just like Extract, but indicates issues with errors, it uses context.Background under the hood. Use ExtractOrErrorContext to specify a context
func ExtractOrError(pdfBytes []byte) ([]PdfPage, error) {
	return ExtractOrErrorContext(context.Background(), pdfBytes)
}

// ExtractOrErrorContext behaves the same as ExtractOrError but lets the caller pass a Context
func ExtractOrErrorContext(ctx context.Context, pdfBytes []byte) ([]PdfPage, error) {
	pages, err := ExtractContext(ctx, pdfBytes)

	if err != nil {
		return pages, err
	}

	if len(pages) == 0 {
		return pages, errors.New("no pages extracted")
	}

	hasContent := false

	for _, p := range pages {
		if p.Content != "" {
			hasContent = true

			break
		}
	}

	if !hasContent {
		return pages, errors.New("no page text extracted")
	}

	return pages, err
}

// ExtractInPopplerTsv Access raw stdout content from Poppler, it uses context.Background under the hood. Use ExtractInPopplerTsvContext to specify a context
func ExtractInPopplerTsv(pdfBytes []byte) ([]PopplerTsvRow, error) {
	return ExtractInPopplerTsvContext(context.Background(), pdfBytes)
}

// ExtractInPopplerTsvContext behaves the same as ExtractInPopplerTsv but lets the caller pass a Context
func ExtractInPopplerTsvContext(ctx context.Context, pdfBytes []byte) ([]PopplerTsvRow, error) {
	params := []string{
		"-tsv",
		"-", // Read from stdin
		"-", // Write to stdout
	}

	cmd := exec.CommandContext(ctx, "pdftotext", params...)
	cmd.Stdin = bytes.NewReader(pdfBytes)

	var out bytes.Buffer

	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, "error executing pdftotext binary")
	}

	var tsvRows []PopplerTsvRow

	tsvT := reflect.TypeFor[PopplerTsvRow]()
	scanner := bufio.NewScanner(strings.NewReader(out.String()))

	scanner.Scan() // Ignore TSV header

	for scanner.Scan() {
		var (
			line   = scanner.Text()
			fields = strings.Fields(line)
		)

		newTsv := PopplerTsvRow{}

		for i := 0; i < tsvT.NumField(); i++ {
			if i >= len(fields) {
				continue
			}

			field := reflect.ValueOf(&newTsv).Elem().Field(i)
			col, err := strconv.Atoi(tsvT.Field(i).Tag.Get("col"))

			if err != nil {
				return nil, errors.Wrap(err, "cannot parse tag as int")
			}

			switch field.Interface().(type) {
			case int:
				newInteger, err := strconv.Atoi(fields[col])

				if err != nil {
					return nil, errors.Wrap(err, "cannot convert value to int")
				}

				field.SetInt(int64(newInteger))
			case float64:
				newFloat, err := strconv.ParseFloat(fields[col], 64)

				if err != nil {
					return nil, errors.Wrap(err, "cannot convert value to float32")
				}

				field.SetFloat(newFloat)
			case string:
				field.SetString(fields[col])
			default:
				return nil, fmt.Errorf("cannot map %s", field.Type().String())
			}
		}

		tsvRows = append(tsvRows, newTsv)
	}

	return tsvRows, nil
}

// CheckPopplerVersion checks the version of the currently-available poppler tool and returns it if suitable, otherwise it will return an error.
func CheckPopplerVersion(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "pdftotext", "-v")

	var out bytes.Buffer

	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", errors.Wrap(err, "error executing binary")
	}

	scanner := bufio.NewScanner(strings.NewReader(out.String()))

	scanner.Scan()
	line := scanner.Text()
	fields := strings.Fields(line)

	if len(fields) < 2 {
		return "", errors.New("no version information extracted")
	}

	fullVersionString := fields[2]

	constraint, err := semver.NewConstraint(popplerVersionConstraint)

	if err != nil {
		return "", fmt.Errorf("cannot parse constraint string \"%s\"", popplerVersionConstraint)
	}

	version, err := semver.NewVersion(fullVersionString)

	if err != nil {
		return "", fmt.Errorf("cannot parse version string \"%s\"", fullVersionString)
	}

	if constraint != nil && version != nil && constraint.Check(version) {
		// poppler is compatible
		return fullVersionString, nil
	}

	return "", fmt.Errorf("incompatible poppler version: require version \"%s\", but found version \"%s\"", constraint.String(), version.String())
}

func init() {
	if _, err := CheckPopplerVersion(context.Background()); err != nil {
		panic(err)
	}
}
