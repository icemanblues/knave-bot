package shakespeare

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

// ComplimentGenerator Generator for Shakespearean Compliments
var ComplimentGenerator = New(Thou, "", [][]string{complimentA, complimentB, complimentC})

// Compliment randomly generates a Shakespearean Compliment
func Compliment() string {
	return ComplimentGenerator.Sentence()
}
