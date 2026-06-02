package utils

import (
	"fmt"
	"strings"
)

var  Voices = map[string]string{
"Achernar": "F",
"Achird": "M",
"Algenib": "M",
"Algieba": "M",
"Alnilam": "M",
"Aoede": "F",
"Autonoe": "F",
"Callirrhoe": "F",
"Charon": "M",
"Despina": "F",
"Enceladus": "M",
"Erinome": "F",
"Fenrir": "M",
"Gacrux": "F",
"Iapetus": "M",
"Kore": "F",
"Laomedeia": "F",
"Leda": "F",
"Orus": "M",
"Puck": "M",
"Pulcherrima": "M",
"Rasalgethi": "M",
"Sadachbia": "M",
"Sadaltager": "M",
"Schedar": "M",
"Sulafat": "F",
"Umbriel": "M",
"Vindemiatrix": "F",
"Zephyr": "F",
"Zubenelgenubi": "M",
}

func PrintVoices(){
	for v, k := range Voices {
		fmt.Printf("%s (%s)\n", v, k)
	}
}

func HasVoice(voice string) bool {
	if len(voice) == 0 {
		return false
	}

	firstLetter := strings.ToUpper(voice[:1])
	rest := strings.ToLower(voice[1:])
	voice = firstLetter + rest

	_, ok := Voices[voice]
	return ok
}
