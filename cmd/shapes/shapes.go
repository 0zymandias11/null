package main

import (
	"math"
)

type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	return r.Height * r.Width
}

func Perimeter(r Rectangle) float64 {
	return 2 * (r.Height + r.Width)
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

// func main() {
// 	// Example usage:
// 	rect := Rectangle{Width: 10, Height: 5}
// 	fmt.Println("Perimeter:", Perimeter(rect))
// 	fmt.Println("Rectangle Area:", rect.Area())
// 	circle := Circle{Radius: 10}
// 	fmt.Println("Circle Area:", circle.Area())
// }
