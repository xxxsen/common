package utils

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(l, r int) string {
	if l > r {
		l = r
	}
	length := secureRandomInt(l, r)
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[secureRandomInt(0, len(letters)-1)]
	}
	return string(b)
}
