package utils

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type ProgressSpinner struct {
	spinner *spinner.Spinner
}

func NewProgressSpinner(message string) *ProgressSpinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" %s", message)
	return &ProgressSpinner{
		spinner: s,
	}
}

func (p *ProgressSpinner) Start() {
	p.spinner.Start()
}

func (p *ProgressSpinner) Stop() {
	p.spinner.Stop()
}

func (p *ProgressSpinner) Success(message string) {
	p.spinner.Stop()
	fmt.Printf("✓ %s\n", message)
}

func (p *ProgressSpinner) Error(message string) {
	p.spinner.Stop()
	fmt.Printf("✗ %s\n", message)
}

func (p *ProgressSpinner) WithDelay(fn func() error, delay time.Duration) error {
	p.Start()

	done := make(chan struct {
		err error
	})

	go func() {
		err := fn()
		done <- struct {
			err error
		}{err}
	}()

	select {
	case <-time.After(delay):
		result := <-done
		if result.err != nil {
			p.Error(result.err.Error())
			return result.err
		}
		return nil
	case result := <-done:
		if delay > 0 {
			time.Sleep(delay)
		}
		if result.err != nil {
			p.Error(result.err.Error())
			return result.err
		}
		return nil
	}
}
