package utils

import "strconv"

func LuhnCheck(number string) bool {
	sum := 0

	for i := 0; i < len(number); i++ {
		num, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}

		if i%2 == len(number)%2 {
			num *= 2
			if num > 9 {
				num -= 9
			}
		}

		sum += num
	}

	return sum%10 == 0
}
