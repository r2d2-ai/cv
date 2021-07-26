package cv

type FPSCounter []int64

func (counter *FPSCounter) FPS() float64 {
	var total int64 = 0
	slice := *counter
	if len(slice) > 1000 {
		*counter = slice[len(*counter)-1000:]
	}

	for _, val := range *counter {
		total += val
	}
	fps := 1000. / (float64(total) / float64(len(*counter)))
	return fps
}
