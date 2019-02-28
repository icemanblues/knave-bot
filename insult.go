package main

import (
	"fmt"
	"math/rand"
	"time"
)

const thou = "Thou"

var column1 = []string{
	"artless",
	"bawdy",
	"beslubbering",
	"bootless",
	"churlish",
	"cockered",
	"clouted",
	"craven",
	"currish",
	"dankish",
	"dissembling",
	"droning",
	"errant",
	"fawning",
	"fobbing",
	"froward",
	"frothy",
	"gleeking",
	"goatish",
	"gorbellied",
	"impertinent",
	"infectious",
	"jarring",
	"loggerheaded",
	"lumpish",
	"mammering",
	"mangled",
	"mewling",
	"paunchy",
	"pribbling",
	"puking",
	"puny",
	"qualling",
	"rank",
	"reeky",
	"roguish",
	"ruttish",
	"saucy",
	"spleeny",
	"spongy",
	"surly",
	"tottering",
	"unmuzzled",
	"vain",
	"venomed",
	"villainous",
	"warped",
	"wayward",
	"weedy",
	"yeasty",
}

var column2 = []string{
	"base-court",
	"bat-fowling",
	"beef-witted",
	"beetle-headed",
	"boil-brained",
	"clapper-clawed",
	"clay-brained",
	"common-kissing",
	"crook-pated",
	"dismal-dreaming",
	"dizzy-eyed",
	"doghearted",
	"dread-bolted",
	"earth-vexing",
	"elf-skinned",
	"fat-kidneyed",
	"fen-sucked",
	"flap-mouthed",
	"fly-bitten",
	"folly-fallen",
	"fool-born",
	"full-gorged",
	"guts-griping",
	"half-faced",
	"hasty-witted",
	"hedge-born",
	"hell-hated",
	"idle-headed",
	"ill-breeding",
	"ill-nurtured",
	"knotty-pated",
	"milk-livered",
	"motley-minded",
	"onion-eyed",
	"plume-plucked",
	"pottle-deep",
	"pox-marked",
	"reeling-ripe",
	"rough-hewn",
	"rude-growing",
	"rump-fed",
	"shard-borne",
	"sheep-biting",
	"spur-galled",
	"swag-bellied",
	"tardy-gaited",
	"tickle-brained",
	"toad-spotted",
	"unchin-snouted",
	"weather-bitten",
}

var column3 = []string{
	"apple-john",
	"baggage",
	"barnacle",
	"bladder",
	"boar-pig",
	"bugbear",
	"bum-bailey",
	"canker-blossom",
	"clack-dish",
	"clotpole",
	"coxcomb",
	"codpiece",
	"death-token",
	"dewberry",
	"flap-dragon",
	"flax-wench",
	"flirt-gill",
	"foot-licker",
	"fustilarian",
	"giglet",
	"gudgeon",
	"haggard",
	"harpy",
	"hedge-pig",
	"horn-beast",
	"hugger-mugger",
	"joithead",
	"lewdster",
	"lout",
	"maggot-pie",
	"malt-worm",
	"mammet",
	"measle",
	"minnow",
	"miscreant",
	"moldwarp",
	"mumble-news",
	"nut-hook",
	"pigeon-egg",
	"pignut",
	"puttock",
	"pumpion",
	"ratsbane",
	"scut",
	"skainsmate",
	"strumpet",
	"varlot",
	"vassal",
	"whey-face",
	"wagtail",
}

func Generate() string {
	rand.Seed(time.Now().Unix())
	r1 := rand.Intn(len(column1))
	r2 := rand.Intn(len(column2))
	r3 := rand.Intn(len(column3))
	return fmt.Sprintf("%s %s %s %s", thou, column1[r1], column2[r2], column3[r3])
}

func main() {
	fmt.Println(Generate())
}
