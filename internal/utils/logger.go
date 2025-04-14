package utils

import (
	"fmt"

	"github.com/fatih/color"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Success(format string, a ...interface{}) {
	green := color.New(color.FgHiGreen).SprintFunc()
	fmt.Printf("%s %s\n", green("✓"), fmt.Sprintf(format, a...))
}

func (l *Logger) Error(format string, a ...interface{}) {
	red := color.New(color.FgHiRed).SprintFunc()
	fmt.Printf("%s %s\n", red("✗"), fmt.Sprintf(format, a...))
}

func (l *Logger) Info(format string, a ...interface{}) {
	blue := color.New(color.FgHiBlue).SprintFunc()
	fmt.Printf("%s %s\n", blue("ℹ"), fmt.Sprintf(format, a...))
}

func (l *Logger) Warning(format string, a ...interface{}) {
	yellow := color.New(color.FgHiYellow).SprintFunc()
	fmt.Printf("%s %s\n", yellow("⚠"), fmt.Sprintf(format, a...))
}

func (l *Logger) Title(format string, a ...interface{}) {
	bold := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	fmt.Println(bold(fmt.Sprintf(format, a...)))
}
