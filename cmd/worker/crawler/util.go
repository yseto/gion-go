package crawler

const DescriptionLength = 255
const TitleLength = 80

func SubstringByBytesLength(s string, blen int) string {
	var output string
	for _, i := range []rune(s) {
		if len(output)+len(string(i)) > blen {
			break
		}
		output += string(i)
	}
	return output
}
