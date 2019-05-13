package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"time"
)

const (
	// CSV url
	cash5Url = "https://www.nclottery.com/Cash5Download"
	// File to read in
	filename = "NCELCash5.csv"

	// Number of times winning combination hit
	cNumCheck = 3
	// Number of times single number hit
	sNumCheck = 580
)

type kv struct {
	Key   string
	Value int
}

func main() {

	// Get CSV from NC Lottery site
	DownloadFile("test.csv", cash5Url)

	// Variables needed
	var wNumbers [][]string
	var cNum = make(map[string]int)
	var sNum = make(map[string]int)
	var cSort []kv
	var sSort []kv

	// Open file, check for err, defer close
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error Opening file: %+v\n", err)
	}
	defer f.Close()

	// Read file and return slice of slice
	wNumbers = ReadFileOut(f)

	// Single winning numbers counted and sorted
	sNum = CountNumbers(wNumbers)
	sSort = SortNumbers(sNum)

	// Combo winning numbers counted and sorted
	cNum = CombineNumbers(wNumbers)
	cSort = SortNumbers(cNum)

	// Print all the numbers (combo & single)
	//PrintAllNumbers(cSort)
	//PrintAllNumbers(sSort)

	// Print only combo numbers that hit more then N times
	PrintMultipleWinningCombos(cSort, cNumCheck)
	PrintMultipleWinningCombos(sSort, sNumCheck)

	myNum := RandomNumGenerator()
	//RandomNumGenerator()

	//fmt.Println(wNumbers)

	fmt.Println("My Number -- ", myNum)

	CheckAgainstWinners(myNum, wNumbers)

}

func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	fmt.Printf("Getting file and saving it as: %+v\n", filepath)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("File download complete")

	return nil
}

func CheckAgainstWinners(mNums []string, wNums [][]string) {
	myNumComb := mNums
	myWNumsComb := wNums

	var myCombo string
	var winNumCombo []string

	for _, v := range myNumComb {
		myCombo += fmt.Sprintf("%s", v)
	}
	//fmt.Println(myCombo)

	for _, v := range myWNumsComb {
		tmp := fmt.Sprintf("%s%s%s%s%s", v[0], v[1], v[2], v[3], v[4])
		winNumCombo = append(winNumCombo, tmp)
	}
	//fmt.Println(myWNumsComb)

	for _, val := range winNumCombo {
		if myCombo == val {
			fmt.Printf("Match -- %s / %s", myNumComb, val)
		}
	}
}

// RandomNumGenerator will randomly generate 5 numbers that
// are between 1 and 43. Calls genCheck to ensure numbers are
// unique (non duplicates).
func RandomNumGenerator() []string {
	var qpNum []string
	var result []string

	rand.Seed(time.Now().UnixNano())

	// loop 5 times to get 5 random numbers,
	// convert to string for handling later,
	// append to slice.
	for i := 1; i <= 5; i++ {
		n := rand.Intn(43) + 1
		//n := rand.Intn(36) + 1
		sn := fmt.Sprintf("%d", n)

		qpNum = append(qpNum, sn)
	}

	// run random generated numbers through
	// genCheck to ensure uniqueness.
	result = genCheck(qpNum)

	// sort unique randomly generated numbers.
	sort.Strings(result)

	//fmt.Println(result)

	return result
}

// genCheck takes slice, calls dupCheck func,
// infinite loop to check length of slice after
// initial dupCheck and will add any missing number
// back to slice. Perform dupCheck again before
// breaking loop.
func genCheck(s []string) []string {
	// perform dupCheck
	nums := dupCheck(s)

	// infinite loop
	for {
		// switch statement check based on length
		// of slice.
		switch len(nums) {
		case 5:
			break
		case 4:
			for i := 4; i < 5; i++ {
				n := rand.Intn(43) + 1
				sn := fmt.Sprintf("%d", n)
				nums = append(nums, sn)
			}
		case 3:
			for i := 3; i < 5; i++ {
				n := rand.Intn(43) + 1
				sn := fmt.Sprintf("%d", n)
				nums = append(nums, sn)
			}
		case 2:
			for i := 2; i < 5; i++ {
				n := rand.Intn(43) + 1
				sn := fmt.Sprintf("%d", n)
				nums = append(nums, sn)
			}
		case 1:
			for i := 1; i < 5; i++ {
				n := rand.Intn(43) + 1
				sn := fmt.Sprintf("%d", n)
				nums = append(nums, sn)
			}
		}

		// Last dupCheck before breaking infinite loop
		nums = dupCheck(nums)
		if len(nums) == 5 {
			break
		}
	}

	return nums
}

// dupCheck takes slice and checks for duplicates,
// if duplicate exists it is removed by creating a
// new slice and appending unique numbers it.
func dupCheck(s []string) []string {
	var encountered = map[string]bool{}
	var result []string

	// range over slice, if in map do nothing,
	// if not, add to map and append to new slice
	// return new slice.
	for _, v := range s {
		if encountered[v] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[v] = true
			// Append to result slice.
			result = append(result, v)
		}
	}

	return result
}

// ReadFileOut takes a csv file and reads in each line,
// adding the individual lotto numbers to slice, then
// adding that slice to a slice that will hold all winning
// combinations.
func ReadFileOut(f *os.File) [][]string {
	var wNums [][]string

	// create Reader for file
	r := csv.NewReader(bufio.NewReader(f))

	// loop over csv and read data in, handle err if one
	// occurs, add column data to slice, add slice to slice
	// of winning number combinations.
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error Reading Line %+v\n", err)
		}

		line = []string{line[1], line[2], line[3], line[4], line[5]}
		sort.Strings(line)
		wNums = append(wNums, line)
	}

	return wNums
}

// CombineNumbers takes winning lotto numbers and combines them
// into a single string and counts the number of times it
// happen to be a winning combination.
func CombineNumbers(n [][]string) map[string]int {
	var m = make(map[string]int)

	// range over slice of slice of winning numbers, make single
	// string, add to map and increment base of number of times it
	// was a winning combination.
	for _, v := range n {
		tmp := fmt.Sprintf("%s%s%s%s%s", v[0], v[1], v[2], v[3], v[4])
		m[tmp] += 1
	}

	return m
}

// CountNumbers takes winning lotto numbers and counts each
// the number of times it was in a winning combination
func CountNumbers(n [][]string) map[string]int {
	var sNum = make(map[string]int)

	// range over slice of slice of winning numbers, add each number
	// to map, increment count based on times it occurred.
	for _, s := range n {
		for _, num := range s {
			sNum[num] += 1
		}
	}

	return sNum
}

// SortNumbers will take map of lotto combination/number of times won
// and sorts the lotto combinations
func SortNumbers(m map[string]int) []kv {
	var ss []kv

	// range over map, append key/value to Slice
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	// sort slice based on Value from map
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return ss
}

// PrintAllNumbers displays all winning lotto combinations that
// is supplied via file.
func PrintAllNumbers(ss []kv) {
	// range over sorted Slice and print out lotto data
	for _, kv := range ss {
		fmt.Printf("%s - %d\n", kv.Key, kv.Value)
	}
}

// PrintMultipleWinningCombos displays winning lotto combinations that
// have been hit multiple times. User specifies the number of multiples
// they are looking for.
func PrintMultipleWinningCombos(ss []kv, multi int) {
	// range over kv struct, check value is more then multi
	// user is specifying.
	for _, val := range ss {
		if val.Value >= multi {
			fmt.Printf("--- %s : %d\n", val.Key, val.Value)
		}
	}
}
