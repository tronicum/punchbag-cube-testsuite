package output

import "fmt"

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func Success(msg string) {
	fmt.Println(ColorGreen + msg + ColorReset)
}

func Error(msg string) {
	fmt.Println(ColorRed + msg + ColorReset)
}

func Warn(msg string) {
	fmt.Println(ColorYellow + msg + ColorReset)
}

func Info(msg string) {
	fmt.Println(ColorBlue + msg + ColorReset)
}
