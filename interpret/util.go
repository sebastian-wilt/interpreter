package interpret

func intPow(left int, right int) int {
	if right == 0 {
		return 1
	}

	if right == 1 {
		return left
	}

	result := left
	for i := 2; i < right; i++ {
		result *= left
	}
	return result
}

func modulo(left int, right int) int {
	rem := left % right
	if rem < 0 {
		rem += right
	}

	return rem
}
