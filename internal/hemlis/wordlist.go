package hemlis

import (
	"fmt"
	"strings"
)

// 256 words from Bytewords / Blockchain Commons (https://github.com/BlockchainCommons/bc-bytewords/blob/c32d8b59690b90f60d26696c955edfecf7e22237/src/bc-bytewords.c
var wordlist = strings.Split("able acid also apex aqua arch atom aunt away axis back bald barn belt beta bias blue body brag brew bulb buzz calm cash cats chef city claw code cola cook cost crux curl cusp cyan dark data days deli dice diet door down draw drop drum dull duty each easy echo edge epic even exam exit eyes fact fair fern figs film fish fizz flap flew flux foxy free frog fuel fund gala game gear gems gift girl glow good gray grim guru gush gyro half hang hard hawk heat help high hill holy hope horn huts iced idea idle inch inky into iris iron item jade jazz join jolt jowl judo jugs jump junk jury keep keno kept keys kick kiln king kite kiwi knob lamb lava lazy leaf legs liar limp lion list logo loud love luau luck lung main many math maze memo menu meow mild mint miss monk nail navy need news next noon note numb obey oboe omit onyx open oval owls paid part peck play plus poem pool pose puff puma purr quad quiz race ramp real redo rich road rock roof ruby ruin runs rust safe saga scar sets silk skew slot soap solo song stub surf swan taco task taxi tent tied time tiny toil tomb toys trip tuna twin ugly undo unit urge user vast very veto vial vibe view visa void vows wall wand warm wasp wave waxy webs what when whiz wolf work yank yawn yell yoga yurt zaps zero zest zinc zone zoom", " ")

func EncodeBytesToWords(data []byte) []string {
	// Convert bytes to a binary string
	var words []string
	for _, b := range data {
		index := int(b)
		words = append(words, wordlist[index])
	}
	return words
}

func DecodeWordsToBytes(encoded []string) ([]byte, error) {
	wordToIndex := make(map[string]int)
	for i, word := range wordlist {
		wordToIndex[word] = i
	}

	var data []byte
	for _, word := range encoded {
		index, ok := wordToIndex[word]
		if !ok {
			return nil, fmt.Errorf("word not found in wordlist: %s", word)
		}
		data = append(data, byte(index))
	}
	return data, nil
}
