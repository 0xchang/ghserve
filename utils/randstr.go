package utils

import "math/rand"

var letters = []rune("1qaz2wsx3edc4rfv5tgb6yhn7ujm8ik9ol0p")

func Random_str() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
