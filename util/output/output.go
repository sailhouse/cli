package output

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	textErr = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff7477"))
)

type Output[T any] struct {
	Messages []string
	Errors   []string
	Table    *Table
	// Used for non-text output
	Data T
}

func WithOutput[T any](f func(cmd *cobra.Command, args []string, output *Output[T])) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		output := NewOutput[T]()
		f(cmd, args, output)

		output.Print()
	}
}

func NewOutput[T any]() *Output[T] {
	return &Output[T]{}
}

func (o *Output[T]) AddMessage(message string) {
	o.Messages = append(o.Messages, message)
}

func (o *Output[T]) AddError(message string, err ...error) {
	if len(err) > 0 {
		message = fmt.Sprintf("%s: %s", message, err[0].Error())
	}

	o.Errors = append(o.Errors, message)
}

func (o *Output[T]) SetData(data T) {
	o.Data = data
}

func (o *Output[T]) SetTable(table *Table) {
	o.Table = table
}

func (o *Output[T]) Print() {
	format := viper.Get("format")

	switch format {
	case "json":
		o.PrintJSON()
	default:
		if o.Table != nil {
			o.Table.Print()
		} else {
			o.PrintText()
		}
	}

	fmt.Println()
}

func (o *Output[T]) PrintErrors() {
	if len(o.Errors) == 1 {
		fmt.Println(textErr.Render(o.Errors[0]))
		return
	}

	fmt.Println(textErr.Render("Errors:"))
	for _, err := range o.Errors {
		fmt.Printf(" - %s\n", textErr.Render(err))
	}
}

func (o *Output[T]) PrintText() {
	if len(o.Errors) > 0 {
		o.PrintErrors()
		return
	}
	for _, message := range o.Messages {
		fmt.Println(message)
	}
}

func (o *Output[T]) PrintJSON() {
	if len(o.Errors) > 0 {
		errBytes, err := json.Marshal(o.Errors)
		if err != nil {
			fmt.Printf("Error marshalling errors: %s\n", err)
		}

		fmt.Printf("{\"errors\": %s}\n", errBytes)
		return
	}

	dataBytes, err := json.Marshal(o.Data)
	if err != nil {
		fmt.Printf("Error marshalling data: %s", err)
	}

	fmt.Printf("%s\n", dataBytes)
}
