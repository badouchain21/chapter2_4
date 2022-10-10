package main
import (
	"fmt"
)

func main()  {
	Diamond(7,7);
}

/**
菱形输出
*/
func Diamond(cmax, rmax int) {
	laststar := 1
	for r := 1; r <= rmax; r++ {
		var star, space, start, end int
		if (r == 1) {
			star = laststar
		} else if (r > 1 && r <= 4) {
			star = laststar + 2
		} else {
			star = laststar - 2
		}
		laststar = star
		space = cmax - star
		start = space/2 + 1
		end = (start + star) - 1
		for c := 1; c <= cmax; c++ {
			if (c < start || c > end) {
				fmt.Print(" ")
			} else {
				fmt.Print("*")
			}
		}
		fmt.Println()
	}
}
