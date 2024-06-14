package util

import (
	"fmt"
	"os"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

type Helper struct {
	cmd *cobra.Command
}

func NewHelper(cmd *cobra.Command) *Helper {
	return &Helper{cmd}
}

func (h *Helper) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

// Print is a convenience method to Print to the defined output, fallback to Stderr if not set.
func (h *Helper) Print(i ...interface{}) {
	h.cmd.Print(i...)
}

// Println is a convenience method to Println to the defined output, fallback to Stderr if not set.
func (h *Helper) Println(i ...interface{}) {
	h.cmd.Println(i...)
}

// Printf is a convenience method to Printf to the defined output, fallback to Stderr if not set.
func (h *Helper) Printf(format string, i ...interface{}) {
	h.cmd.Printf(format, i...)
}

// PrintErr is a convenience method to Print to the defined Err output, fallback to Stderr if not set.
func (h *Helper) PrintErr(i ...interface{}) {
	h.cmd.PrintErr(i...)
}

// PrintErrln is a convenience method to Println to the defined Err output, fallback to Stderr if not set.
func (h *Helper) PrintErrln(i ...interface{}) {
	h.cmd.PrintErrln(i...)
}

// PrintErrf is a convenience method to Printf to the defined Err output, fallback to Stderr if not set.
func (h *Helper) PrintErrf(format string, i ...interface{}) {
	h.cmd.PrintErrf(format, i...)
}

func (h *Helper) Status(i ...interface{}) {
	h.cmd.Printf("[%s] ", h.cmd.CommandPath())
	h.cmd.Println(i...)
}

func (h *Helper) Statusf(format string, i ...interface{}) {
	h.Status(fmt.Sprintf(format, i...))
}

func (h *Helper) Fatal(args ...interface{}) {
	h.PrintErrln(args...)
	os.Exit(1)
}

func (h *Helper) Fatalf(format string, args ...interface{}) {
	h.Fatal(fmt.Sprintf(format, args...))
}

func (h *Helper) Table() *uitable.Table {
	table := uitable.New()
	table.MaxColWidth = 50
	table.Separator = "  "
	return table
}
