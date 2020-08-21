package main

import (
	"fmt"
	"log"
	"reflect"
	"sort"
)

//Seat Structure
type Seat struct {
	index    int
	row      int
	column   int
	category string
}

//less function for sorter Interface
type lessFunc func(p1, p2 *Seat) bool

//Sorter Interface  -- Sorts the Seat instance
type multiSorter struct {
	seats []Seat
	less  []lessFunc
}

// multiSorter implementation for Sorting
func (ms *multiSorter) Sort(seats []Seat) {
	ms.seats = seats
	sort.Sort(ms)
}

//  sort.Interface implementation
func (ms *multiSorter) Len() int {
	return len(ms.seats)
}

//  sort.Interface implementation
func (ms *multiSorter) Swap(i, j int) {
	ms.seats[i], ms.seats[j] = ms.seats[j], ms.seats[i]
}

//  sort.Interface implementation
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.seats[i], &ms.seats[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}

	}
	return ms.less[k](p, q)
}

//OrderBy Implementation - MultiSorting functionality
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Identify Windows Seats
// param bayIndex - Identifies the bay index
// param cornerBayFlag - true/false identifies first/last bay flag
// param bayDimensions - size of the given bay
// returns list of Window Seats as []Seat
func identifyWindowSeats(bayIndex int, cornerBayFlag bool, bayDimensions []int) []Seat {

	flightRows := bayDimensions[1]
	flightColumns := bayDimensions[0]
	var seatList []Seat

	if cornerBayFlag {
		for i := 1; i <= flightRows; i++ {
			seatList = append(seatList, Seat{bayIndex + 1, i, 1, "B"})
		}
	} else {
		for i := 1; i <= flightRows; i++ {
			seatList = append(seatList, Seat{bayIndex + 1, i, flightColumns, "B"})
		}
	}
	return seatList

}

// Identify Aisle Seats
// param bayIndex - Identifies the bay index
// param cornerBayFlag - F/B/C identifies Front,Back &Centre
// param bayDimensions - size of the given bay
// returns list of Aisle Seats as []Seat
func identifyAisleSeats(bayIndex int, cornerBayFlag string, bayDimensions []int) []Seat {

	flightRows := bayDimensions[1]
	flightColumns := bayDimensions[0]
	var seatList []Seat

	if cornerBayFlag == "F" {
		for i := 1; i <= flightRows; i++ {
			seatList = append(seatList, Seat{bayIndex + 1, i, flightColumns, "A"})
		}
	} else if cornerBayFlag == "B" {
		for i := 1; i <= flightRows; i++ {
			seatList = append(seatList, Seat{bayIndex + 1, i, 1, "A"})
		}
	} else {
		for i := 1; i <= flightRows; i++ {
			seatList = append(seatList, Seat{bayIndex + 1, i, 1, "A"})
			seatList = append(seatList, Seat{bayIndex + 1, i, flightColumns, "A"})
		}

	}

	return seatList
}

// Identify Middle Seats
// param bayIndex - Identifies the bay index
// param bayDimensions - size of the given bay
// returns list of Middle Seats as []Seat
func identifyMiddleSeats(bayIndex int, bayDimensions []int) []Seat {

	flightRows := bayDimensions[1]
	flightColumns := bayDimensions[0]
	var seatList []Seat

	for i := 1; i <= flightRows; i++ {
		for j := 1; j <= flightColumns; j++ {
			if j != 1 && j != flightColumns {
				seatList = append(seatList, Seat{bayIndex + 1, i, j, "M"})
			}
		}
	}
	return seatList
}

// Identify Seat Map : Identifies Window, Aisle & Middle Seats
// inputDimensions : Dimension provided for each bay
// returns Seat Mapping as []Seat
func identifySeatMap(inputDimensions [][]int) []Seat {

	var finalSeatList []Seat
	for i, v := range inputDimensions {
		if i == 0 {
			windowSeats := identifyWindowSeats(i, true, v)
			finalSeatList = append(finalSeatList, windowSeats...)

			aisleSeats := identifyAisleSeats(i, "F", v)
			finalSeatList = append(finalSeatList, aisleSeats...)

		} else if i == len(inputDimensions)-1 {
			windowSeats := identifyWindowSeats(i, false, v)
			finalSeatList = append(finalSeatList, windowSeats...)

			aisleSeats := identifyAisleSeats(i, "B", v)
			finalSeatList = append(finalSeatList, aisleSeats...)

		} else {
			aisleSeats := identifyAisleSeats(i, "C", v)
			finalSeatList = append(finalSeatList, aisleSeats...)
		}

		middleSeats := identifyMiddleSeats(i, v)
		finalSeatList = append(finalSeatList, middleSeats...)

	}

	return finalSeatList

}

// Lambdas to sort based on index, row, column category for a list of Seats
// input seatList identified with Windows, Aisle, Center identifers
func sortSeatMap(finalSeatList []Seat) {

	// Sort by Bay Index
	index := func(c1, c2 *Seat) bool {
		return c1.index < c2.index
	}

	// Sort by Row
	row := func(c1, c2 *Seat) bool {
		return c1.row < c2.row
	}

	// Sort by Column
	column := func(c1, c2 *Seat) bool {
		return c1.column < c2.column
	}

	// Sort by Category
	category := func(c1, c2 *Seat) bool {
		return c1.category < c2.category
	}

	// Final Order as part the requirement
	// Aisle First - A
	// Window Second - B
	// Center - C
	OrderedBy(category, row, index, column).Sort(finalSeatList)

}

// Print Seat Mapping as per the allocations
// inputDimensions : Dimension provided for each bay
// finalSeatList mapped, sorted and rank assigned
// passengerQueueLenght
func printSeatMap(inputDimensions [][]int, finalSeatList []Seat, passengerQueueLength int) {

	totalColumns := 0
	totalRows := -1

	for _, v := range inputDimensions {
		totalColumns += v[0]
		if totalRows < v[1] {
			totalRows = v[1]
		}
	}

	finalArray := make([][]int, totalRows)
	for i := range finalArray {
		finalArray[i] = make([]int, totalColumns)
	}

	for i, seat := range finalSeatList {
		if i == passengerQueueLength {
			break
		}

		x := getIndexValue(inputDimensions, seat)
		//fmt.Println(seat , x , seat.row, i+1)
		finalArray[seat.row-1][x-1] = i + 1
	}

	for _, x := range finalArray {
		for _, y := range x {
			fmt.Printf("|%02d|", y)
		}
		fmt.Print("\n")
	}

}

// Function to get the concatenated index value after joining all the arrays
// inputDimensions : Dimension provided for each bay
func getIndexValue(inputDimensions [][]int, seat Seat) int {

	indexX := 0

	for i, v := range inputDimensions {
		if i+1 < seat.index {
			indexX += v[0]
		} else {
			break
		}

	}
	indexX += seat.column

	return indexX
}

// Read Input Parameters
// returns input passenger Array, passengerQueueLength
func readInput() ([][]int, int) {
	inputArray := [][]int{{3, 2}, {4, 3}, {2, 3}, {3, 4}}
	passengerQueue := 30
	return inputArray, passengerQueue
}

// Main Controller function

func process(value [][]int, passengerQueueLength int) {

	finalSeatList := identifySeatMap(value)
	sortSeatMap(finalSeatList)
	printSeatMap(value, finalSeatList, passengerQueueLength)

}

// Main Function to Trigger the program
//*******************************
//****THE PROGRAM STARTS HERE****
//*******************************
func main() {

	//fmt.Println("start")

	executeTests := false
	if executeTests {
		//fmt.Println("Executing Tests")
		testIdentifyWindowSeats()
		testIdentifyAisleSeats()
		testIdentifyMiddleSeats()
		testIdentifySeatMap()
		testSortSeatMap()
		testPrintSeatMap()
		testProcess()
		testGetIndexValue()
		//fmt.Println("Test Execution Complete")
	} else {
		inputArray, passengerQueueLength := readInput()
		//validate()
		process(inputArray, passengerQueueLength)

	}

	//fmt.Println("end")
}

// Unit test
func testIdentifySeatMap() {
	inputArray1 := [][]int{{2, 2}, {3, 3}}
	actualOutput1 := identifySeatMap(inputArray1)
	expectedOutput1 := []Seat{
		{1, 1, 1, "B"},
		{1, 2, 1, "B"},
		{1, 1, 2, "A"},
		{1, 2, 2, "A"},
		{2, 1, 3, "B"},
		{2, 2, 3, "B"},
		{2, 3, 3, "B"},
		{2, 1, 1, "A"},
		{2, 2, 1, "A"},
		{2, 3, 1, "A"},
		{2, 1, 2, "M"},
		{2, 2, 2, "M"},
		{2, 3, 2, "M"},
	}
	if !reflect.DeepEqual(actualOutput1, expectedOutput1) {
		log.Panic("Error executing tests")
	}
}

// Integration test
func testProcess() {
	//TODO: SORT SEATMAP
}

// test case functions
func testGetIndexValue() {
	//TODO: GET INDEX VALUE
}

// test case functions
func testSortSeatMap() {
	//TODO: SORT SEATMAP
}

// test case functions
func testPrintSeatMap() {
	//TODO: PRINT SEATMAP
}

// test case functions
func testIdentifyWindowSeats() {

	input1 := []int{1, 1}
	actualOutput1 := identifyWindowSeats(0, false, input1)
	expectedOutput1 := []Seat{{1, 1, 1, "B"}}
	if !reflect.DeepEqual(actualOutput1, expectedOutput1) {
		log.Panic("Error executing tests")
	}

	fmt.Println("***")
	input2 := []int{3, 2}
	actualOutput2 := identifyWindowSeats(0, false, input2)
	expectedOutput2 := []Seat{{1, 1, 3, "B"}, {1, 2, 3, "B"}}
	if !reflect.DeepEqual(actualOutput2, expectedOutput2) {
		log.Panic("Error executing tests")
	}

	input3 := []int{1, 1}
	actualOutput3 := identifyWindowSeats(0, true, input3)
	expectedOutput3 := []Seat{{1, 1, 1, "B"}}
	if !reflect.DeepEqual(actualOutput3, expectedOutput3) {
		log.Panic("Error executing tests")
	}

	input4 := []int{3, 2}
	actualOutput4 := identifyWindowSeats(0, true, input4)
	expectedOutput4 := []Seat{{1, 1, 1, "B"}, {1, 2, 1, "B"}}
	if !reflect.DeepEqual(actualOutput4, expectedOutput4) {
		log.Panic("Error executing tests")
	}

}

// test case functions
func testIdentifyAisleSeats() {

	input1 := []int{1, 1}
	actualOutput1 := identifyAisleSeats(0, "F", input1)
	expectedOutput1 := []Seat{{1, 1, 1, "A"}}
	if !reflect.DeepEqual(actualOutput1, expectedOutput1) {
		log.Panic("Error executing tests 1")
	}

	fmt.Println("***")
	input2 := []int{3, 2}
	actualOutput2 := identifyAisleSeats(0, "F", input2)
	expectedOutput2 := []Seat{{1, 1, 3, "A"}, {1, 2, 3, "A"}}
	if !reflect.DeepEqual(actualOutput2, expectedOutput2) {
		log.Panic("Error executing tests 2")
	}

	input3 := []int{1, 1}
	actualOutput3 := identifyAisleSeats(0, "B", input3)
	expectedOutput3 := []Seat{{1, 1, 1, "A"}}
	if !reflect.DeepEqual(actualOutput3, expectedOutput3) {
		log.Panic("Error executing tests 3")
	}

	input4 := []int{3, 2}
	actualOutput4 := identifyAisleSeats(0, "B", input4)
	expectedOutput4 := []Seat{{1, 1, 1, "A"}, {1, 2, 1, "A"}}
	if !reflect.DeepEqual(actualOutput4, expectedOutput4) {
		log.Panic("Error executing tests 4")
	}

	input5 := []int{3, 4}
	actualOutput5 := identifyAisleSeats(0, "B", input5)
	expectedOutput5 := []Seat{{1, 1, 1, "A"}, {1, 2, 1, "A"}, {1, 3, 1, "A"}, {1, 4, 1, "A"}}
	if !reflect.DeepEqual(actualOutput5, expectedOutput5) {
		log.Panic("Error executing tests 4")
	}

	input6 := []int{4, 3}
	actualOutput6 := identifyAisleSeats(0, "C", input6)
	expectedOutput6 := []Seat{{1, 1, 1, "A"}, {1, 1, 4, "A"}, {1, 2, 1, "A"}, {1, 2, 4, "A"}, {1, 3, 1, "A"}, {1, 3, 4, "A"}}
	if !reflect.DeepEqual(actualOutput6, expectedOutput6) {
		log.Panic("Error executing tests 4")
	}

}

//  Unit test case functions
func testIdentifyMiddleSeats() {

	input1 := []int{2, 2}
	actualOutput1 := identifyMiddleSeats(0, input1)

	if len(actualOutput1) != 0 {
		log.Panic("Error executing tests 1")
	}

	input2 := []int{3, 2}
	actualOutput2 := identifyMiddleSeats(1, input2)
	expectedOutput2 := []Seat{{2, 1, 2, "M"}, {2, 2, 2, "M"}}
	if !reflect.DeepEqual(actualOutput2, expectedOutput2) {
		log.Panic("Error executing tests 4")
	}

}

