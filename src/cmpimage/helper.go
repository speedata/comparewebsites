package cmpimage

func longestCommonPrefix(str1, str2 string) string {

	r1 := []rune(str1)
	r2 := []rune(str2)

	l1 := len(r1)
	l2 := len(r2)
	if l1 == 0 || l2 == 0 {
		return ""
	}

	if l2 < l1 {
		r1, r2 = r2, r1
		l1 = l2
	}
	lastSlash := 0
	for i := 0; i < l1; i++ {
		if r1[i] != r2[i] {
			if l1 > lastSlash {
				lastSlash++
			}
			return string(r1[:lastSlash])
		}
		if r1[i] == '/' {
			lastSlash = i
		}
	}
	return string(r1)
}
