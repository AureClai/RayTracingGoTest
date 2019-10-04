# RayTracingGoTest
Testing Golang implementation of C++ RayTracing Engine from 
<https://github.com/RayTracing/raytracing.github.io>

For the moment, only a sample scene can be rendered.

## Prerequisites

Go have to be installed on the computer.

## Installation

Clone the repository from GitHub

```Shell
git clone https://github.com/AureClai/RayTracingGoTest
```

## Use 

1. Edit `main.go` to modify the parameters of the output image and sampling (lines 23,26,29) :
``` Go
/*
 	WIDTH is the number of pixel on the X-Axis > 200 recommended
 	HEIGHT is the number of pixels on the Y-Axis > 200 recommended
 	SAMPLES is the number of random rays used to estimate the color of one pixel :
 		- > 10 for first test
 		- > 100 for pretty good render
 		- > 1000 for very 
good render - VERY lONG */
// WIDTH unexported
const WIDTH int = 400

// HEIGHT unexported
const HEIGHT int = 400

// SAMPLES unexported
const SAMPLES int = 100
```

2. Have fun editting the wall colors (lines 34+) - Bad colors values raise a panic
```Go
// WALL COLORS (float64 beetween 0.0 and 1.0)
// initial red(0.65,0.05,0.5) for left and green(0.12,0.45,0.15) for right
// Left Wall
const LeftWallR float64 = 0.65
const LeftWallG float64 = 0.05
const LeftWallB float64 = 0.05

//Right Wall
const RightWallR float64 = 0.45
const RightWallG float64 = 0.45
const RightWallB float64 = 0.45
```
3. Navigate through the system to the folder of the project
4. Run the program with the command
```Shell
go run main.go
```
5. After the computation, the result is next to `main.go` and named `outputImage.ppm`

You can also build the program via `go build main.go`and launch the binary.

## Some ideas for the future

* Make the scene entirely customizable from external files (JSON probably)
* Implement `.obj` support
* Make the code idomatic
* Optimization of the code + Parallelization
* Implement more complex materials