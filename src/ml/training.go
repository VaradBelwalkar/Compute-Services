package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/linearreg"
)

func main() {
	// Load the data
	x := mat.NewDense(3, 2, []float64{1, 2, 3, 4, 5, 6})
	y := mat.NewVecDense(3, []float64{2, 4, 6})

	// Create a linear regression model
	model := linearreg.New(2, false)

	// Train the model using the data
	_, err := model.Fit(x, y)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the model coefficients
	fmt.Println("Coefficients:", model.Coefficients())
	// Output: Coefficients: [1 2]
}
