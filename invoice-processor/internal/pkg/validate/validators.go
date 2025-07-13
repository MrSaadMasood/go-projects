package validate

import (
	"fmt"
	"time"
)

type FlagsOptions struct {
	DryRun bool
	From   string
	To     string
	Path   string
}

func Date(str string) (*time.Time, error) {
	date, err := time.Parse("2006-01-02", str)
	if err != nil {
		return nil, fmt.Errorf("Date Validation Failed %w", err)
	}
	return &date, nil
}

type ValidatedFlags struct {
	From   *time.Time
	To     *time.Time
	DryRun bool
	Path   string
}

func Flags(flags FlagsOptions) (*ValidatedFlags, error) {

	var from *time.Time
	var to *time.Time
	var err error

	if flags.From != "" {
		from, err = Date(flags.From)
		if err != nil {
			return nil, err
		}

	}

	if flags.To != "" {
		to, err = Date(flags.To)
		if err != nil {
			return nil, err
		}
	}

	return &ValidatedFlags{From: from, To: to, DryRun: flags.DryRun, Path: flags.Path}, nil
}
