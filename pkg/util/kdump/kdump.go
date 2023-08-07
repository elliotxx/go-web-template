// Package kdump like fmt.Println but more pretty and beautiful print Go values.
package kdump

import (
	"io"
	"os"

	"github.com/gookit/goutil/dump"
)

var std = (KDumper)(*dump.Std())

func New() KDumper {
	return (KDumper)(*dump.NewDumper(os.Stdout, 3))
}

// V like fmt.Println, but the output is clearer and more beautiful
func V(vs ...any) {
	std.Dump(vs...)
}

// P like fmt.Println, but the output is clearer and more beautiful
func P(vs ...any) {
	std.Print(vs...)
}

// Print like fmt.Println, but the output is clearer and more beautiful
func Print(vs ...any) {
	std.Print(vs...)
}

// Println like fmt.Println, but the output is clearer and more beautiful
func Println(vs ...any) {
	std.Println(vs...)
}

// Fprint like fmt.Println, but the output is clearer and more beautiful
func Fprint(w io.Writer, vs ...any) {
	std.Fprint(w, vs...)
}

// Format like fmt.Println, but the output is clearer and more beautiful
func Format(vs ...any) string {
	return std.Format(vs...)
}

// Custom method outside the original dump package
// FormatN like fmt.Println, but the output is clearer and no color
func FormatN(vs ...any) string {
	return std.FormatN(vs...)
}
