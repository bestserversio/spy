package utils

func GetRegion(region int) string {
	ret := ""

	switch region {
	case 0:
		ret = "NA"
	case 1:
		ret = "NA"
	case 2:
		ret = "SA"
	case 3:
		ret = "EU"
	case 4:
		ret = "AS"
	case 5:
		ret = "AU"
	case 6:
		ret = "ME"
	case 7:
		ret = "AF"
	}

	return ret
}
