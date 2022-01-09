package converter

func pixelToAscii(pixelValue float32) byte {
	switch {
	case pixelValue < 0.1:
		return '.'
	case pixelValue < 0.2:
		return ','
	case pixelValue < 0.3:
		return ';'
	case pixelValue < 0.4:
		return '!'
	case pixelValue < 0.5:
		return 'v'
	case pixelValue < 0.6:
		return 'l'
	case pixelValue < 0.7:
		return 'L'
	case pixelValue < 0.8:
		return 'F'
	case pixelValue < 0.9:
		return 'E'
	default:
		return '$'
	}
}
