package terminfo

import "fmt"

func GoUpBy(amount int) string {
	return fmt.Sprintf("\033[%dA", amount)
}
