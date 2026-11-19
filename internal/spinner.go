package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

func StartSpinner() (*spinner.Spinner, error) {
	s := spinner.New(
		spinner.CharSets[14],
		100*time.Millisecond,
		spinner.WithWriterFile(os.Stderr),
	)
	s.Suffix = " Loading Terraform provider schemas..."
	if err := s.Color("magenta", "bold"); err != nil {
		return nil, fmt.Errorf("failed to set spinner color: %w", err)
	}

	s.Start()
	return s, nil
}

func StopSpinner(s *spinner.Spinner) {
	if s != nil {
		s.Stop()
	}
}
