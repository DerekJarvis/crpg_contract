package chainTypes

import "math"

type ColorAttributes struct {
	AlphaK  int //0
	Red     int //1
	Green   int //2
	Blue    int //3
	Magenta int //4
	Yellow  int //5
	Cyan    int //6
}

func ConvertColorsToArray(attributes *ColorAttributes) []*int {
	valuesArray := []*int{&attributes.AlphaK, &attributes.Red, &attributes.Green, &attributes.Blue, &attributes.Magenta, &attributes.Yellow, &attributes.Cyan}

	return valuesArray
}

func GetHighestStat(attributes *ColorAttributes) *int {
	valuesArray := ConvertColorsToArray(attributes)

	currentMax := valuesArray[0]
	for _, v := range valuesArray[1:] {
		if *v > *currentMax {
			currentMax = v
		}
	}

	return currentMax
}

/// Produces the net outcome between two sets of properties
func SmashColors(attacker *ColorAttributes, defender *ColorAttributes) int {
	// attackerColors := ConvertColorsToArray(attacker)
	// defenderColors := ConvertColorsToArray(attacker)

	// Doing some future proofing in case we add more colors and the array indexing gets broken
	// Unfortunately, there's no assert in golang, so we have to panic at the disco
	// if len(attackerColors) > 7 {
	// 	panic("Color attributes no longer has a length of 7. Things are getting broken.")
	// }

	colorDistanceFloat := math.Sqrt(
		math.Pow(float64(attacker.AlphaK-defender.AlphaK), 2) +
			math.Pow(float64(attacker.Red - -defender.Cyan), 2) +
			math.Pow(float64(attacker.Green - -defender.Magenta), 2) +
			math.Pow(float64(attacker.Blue - -defender.Yellow), 2))

	colorDistanceInt := int(math.Floor(colorDistanceFloat))

	return colorDistanceInt

	// Calculate the distance between the remaining vectors, setting attacker as positive and defender as negative

}
