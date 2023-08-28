package ut

import "math/rand"

var challengeWords = []string{
	"amplify",
	"healer",
	"hexagon",
	"hoist",
	"ambivalent",
}

func GetChallengeWord() string {
	index := rand.Intn(len(challengeWords))
	return challengeWords[index]
}
