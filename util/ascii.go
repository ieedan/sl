package util

const (
	STRIKE_THROUGH_START = "\033[9m"
	STRIKE_THROUGH_END   = "\033[0m"
)

func StrikeThrough(text string) string {
	return STRIKE_THROUGH_START + text + STRIKE_THROUGH_END
}
