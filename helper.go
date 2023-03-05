package main

import (
	"math"
	"strings"
)

// funky function for centering text because i was bored
func centerPaddedString(str string, padRune rune, width int) string {

	var halfwidth int

	centeredLines := []string{}
	lines := strings.Split(str, "\n")

	for _, line := range lines {
		// shitty way to mke sure the width is an even number despite odd # f characters in input text
		if len(line)%2 != 0 {
			line = line + " "
		}
		if width < len(line) {
			width = len(line)
		}
		halfwidth = int(math.Ceil(float64(width/2)-float64(len(line)/2))) + 2 // the two is just to make minimum spacing

		tempString := ""
		for {
			tempString = tempString + " "
			if len(tempString) > halfwidth {
				tempString = tempString[0:halfwidth]
				break
			}

		}

		centeredLines = append(centeredLines, string(padRune)+tempString+line+tempString+string(padRune)+"\n")
	}
	maxCenteredLineLength := maxStringLength(centeredLines)
	centeredText := strings.Join(centeredLines, "")
	border := strings.Repeat(string(padRune), maxCenteredLineLength-1) + "\n"
	edge := string(padRune) + strings.Repeat(" ", maxCenteredLineLength-3) + string(padRune) + "\n"
	formattedText := border + edge + centeredText + edge + border
	return formattedText

}

func maxStringLength(lines []string) int {
	var maxStringLength int = 0
	for i, line := range lines {
		if i == 0 || len(line) > maxStringLength {
			maxStringLength = len(line)
		}
	}
	return maxStringLength
}
