package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var length, width, height, palletWeight, volumetricWeight float64

func calculateVolumetricWeight() float64 {
	return (length * width * height) / 6000
}

func getPalletDetails() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("UK or EU pallet: ")
	scanner.Scan()

	pt := strings.ToLower(scanner.Text())
	switch pt {
	case "uk":
		length, width = 120, 100
	case "eu":
		length, width = 120, 80
	default:
		fmt.Println("Invalid pallet type. Please enter a valid pallet type")
		return
	}

	fmt.Print("Enter the height of the pallet (in cm): ")
	scanner.Scan()
	heightInput := scanner.Text()
	if heightValue, err := strconv.ParseFloat(heightInput, 64); err == nil {
		height = heightValue
	} else {
		fmt.Println("Invalid height value. Please enter a valid number")
		return
	}

	fmt.Print("Enter the weight of the pallet (in cm): ")
	scanner.Scan()
	weightInput := scanner.Text()
	if weightValue, err := strconv.ParseFloat(weightInput, 64); err == nil {
		palletWeight = weightValue
	} else {
		fmt.Println("Invalid weight value. Please enter a valid number")
		return
	}

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	addMore := "yes"
	cw := 0.0
	for addMore == "yes" {
		getPalletDetails()
		volumetricWeight = calculateVolumetricWeight()
		cw += volumetricWeight
		fmt.Printf("Calculated volumetric weight: %.2f kg\n", volumetricWeight)

		fmt.Print("Do you want to add more pallets? (yes/no): ")
		scanner.Scan()
		addMore = strings.ToLower(scanner.Text())
	}

	fmt.Printf("Total Chargeable Weight: %.2f kg\n", cw)
}
