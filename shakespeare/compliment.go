package shakespeare

import (
	"fmt"
	"math/rand"
	"time"
)

var complimentA = []string{
	"rare",
	"sweet",
	"fruitful",
	"brave",
	"sugared",
	"flowering precious",
	"gallant",
	"delicate",
	"celestial",
}

var complimentB = []string{
	"honey-tongued",
	"well-wishing",
	"fair-faced",
	"best-tempered",
	"tender-hearted",
	"tiger-booted",
	"smooth-faced",
	"thunder-darting",
	"sweet-suggesting",
	"young-eyed",
}

var complimentC = []string{
	"smilet",
	"toast",
	"cukoo-bud",
	"nose-herb",
	"wafer-cake",
	"pigeon-egg",
	"welsh cheese",
	"song",
	"true-penny",
	"valentine",
}

// Compliment randomly generates a Shakespearean Compliment
func Compliment() string {
	rand.Seed(time.Now().Unix())
	r1 := rand.Intn(len(complimentA))
	r2 := rand.Intn(len(complimentB))
	r3 := rand.Intn(len(complimentC))
	return fmt.Sprintf("%s %s %s %s", Thou, complimentA[r1], complimentB[r2], complimentC[r3])
}
