package commands

// ANSI escape codes to be used in the color output of the board view
const (
	TERM_RESET = "\033[0m"

	TERM_BG_GRAY  = "\033[48;2;127;127;127m"
	TERM_BG_WHITE = "\033[107m"

	TERM_FG_GRAY      = "\033[38;2;127;127;127m"
	TERM_FG_LIGHTGRAY = "\033[38;2;200;200;200m"
	TERM_FG_FOOD      = "\033[38;2;255;92;117m"
	TERM_FG_RGB       = "\033[38;2;%d;%d;%dm"
)
