package utils

func GetRegion(region int) string {
	ret := ""

	switch region {
	case 0:
		ret = "US"

		break
	case 1:
		ret = "US"

		break

	case 2:
		ret = "SA"

		break

	case 3:
		ret = "EU"

		break

	case 4:
		ret = "AS"

		break

	case 5:
		ret = "AU"

		break

	case 6:
		ret = "ME"

		break

	case 7:
		ret = "AF"

		break
	}

	return ret
}
