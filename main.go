package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/stack"

	//"io"
)

var helpText string = `
this is a simple stack-based calculator program

the following are some currently supported commands:

    q    quits the program
    l    view stack
    +    add last 2 numbers in the stack
    -    subtract the last number from the second-to-last number
    *    multiply the last 2 numbers in the stack
    /    divide the last number from the second-to-last number
    h    print this helptext

enter any number to push it to the stack

multiple commands can be entered on a single line by seperating them with a space

any command line arguments will be parsed as if they were entered interactively
`

var (
	progRunning  = true
	valStack     = stack.New()
	storedVals   = make(map[interface{}]interface{})
	bracketDepth = 0
)

func main() {

	defer showStack(valStack)

	boolFlags, stringFlags := setFlags()

	if len(flag.Args()) >= 1 {
		for _, uinput := range flag.Args() {
			processInput(uinput)
		}
		if !*boolFlags["no-arg-auto-exit"] {
			progRunning = false
		}
	}

	var reader *bufio.Reader
	var readingFile bool

	if *stringFlags["file-path"] != "" {
		inputFile, err := os.Open(*stringFlags["file-path"])
		ec(err)
		reader = bufio.NewReader(inputFile)
		readingFile = true
	} else {
		reader = bufio.NewReader(os.Stdin)
		readingFile = false
	}

	for progRunning {
		fmt.Print("-> ")

		uinputFull, err := reader.ReadString('\n')
		if err != nil {
			if readingFile == true {
				reader = bufio.NewReader(os.Stdin)
			} else {
				panic(err)
			}
		}

		uinputFull = strings.Replace(uinputFull, "\n", "", -1)
		//uinput = strings.Replace(uinput," ","",-1)

		uinputList := strings.SplitN(uinputFull, " ", -1)

		for _, uinput := range uinputList {
			halt := processInput(uinput)
			if halt {
				break
			}
		}
	}
}

func setFlags() (map[string]*bool, map[string]*string) {
	boolFlags := make(map[string]*bool)
	boolFlags["no-arg-auto-exit"] = flag.Bool("no-arg-auto-exit", true, "makes the program not exit automaticaly when run with command-line arguments")
	stringFlags := make(map[string]*string)
	stringFlags["file-path"] = flag.String("f", "", "parse instructions from file")

	flag.Parse()

	return boolFlags, stringFlags
}

func processInput(uinput string) (halt bool) {
	//defer meditate()

	numInput, inpIsNum := strconv.ParseFloat(uinput, 64)

	if len(uinput) == 0 {
		return false
	} else if bracketDepth > 0 {
		val1 := valStack.Pop() //brackets are used to store functions
		valStack.Push(strings.Join([]string{val1.(string), uinput}, " "))
		if uinput == "]" {
			bracketDepth--
		} else if uinput == "[" {
			bracketDepth++
		}
	} else if inpIsNum == nil {
		valStack.Push(numInput)
		fmt.Println(numInput, "pushed")
	} else if uinput[:1] == "\"" && uinput[len(uinput)-1:] == "\"" {
		valStack.Push(uinput)
		fmt.Println(uinput, "pushed")
	} else {
		switch uinput {
		case "q":
			progRunning = false
		case "+":
			val1, val2, err := pop2Vals(valStack)
			if ifErrStackWarn(err, 2) {
				break
			}

			valStack.Push(val1 + val2)
			fmt.Println("sum is", valStack.Peek())
		case "-":
			val1, val2, err := pop2Vals(valStack)
			if ifErrStackWarn(err, 2) {
				break
			}

			valStack.Push(val2 - val1)
			fmt.Println("diff is", valStack.Peek())
		case "*":
			val1, val2, err := pop2Vals(valStack)
			if ifErrStackWarn(err, 2) {
				break
			}

			valStack.Push(val1 * val2)
			fmt.Println("product is", valStack.Peek())
		case "/":
			val1, val2, err := pop2Vals(valStack)
			if ifErrStackWarn(err, 2) {
				break
			}

			valStack.Push(val2 / val1)
			fmt.Println("quotient is", valStack.Peek())
		case "|": // mirror/swap function
			val1 := valStack.Pop()
			val2 := valStack.Pop()

			valStack.Push(val1)
			valStack.Push(val2)

			fmt.Println("values", val2, "and", val1, "swapped")
		case "$": // store function
			val1 := valStack.Pop() //index
			val2 := valStack.Pop() //value

			storedVals[val1] = val2
			fmt.Println("value", val2, "stored under index", val1)
		case "=":
			val1 := valStack.Peek()
			if storedVals[val1] != nil {
				valStack.Push(storedVals[val1])

				fmt.Println("value", valStack.Peek(), "retrived from index", val1)
			} else {
				fmt.Println("index", val1, "does not have a value assigned")
			}
		case "~": //removes items from stack
			val1 := valStack.Pop()

			fmt.Println("value", val1, "removed from stack")
		case "_": //duplicates top item in stack
			val1 := valStack.Peek()
			valStack.Push(val1)
		case "[": //lets you store functions
			bracketDepth = 1
			valStack.Push(uinput)
		case "#": //lets you run stored functions
			instructionsRaw, _ := valStack.Pop().(string)
			instructionsSlice := strings.Split(instructionsRaw, " ")
			instructionsStripped := instructionsSlice[1 : len(instructionsSlice)-1]
			

			for _, instruction := range instructionsStripped {
				halt := processInput(instruction)
				if halt {
					break
				}
			}
		case "?":
			val1, isNum := valStack.Peek().(float64) 
			if val1 == 0 && isNum {
				return true
			}
		case "l":
			showStack(valStack)
		case "h":
			fmt.Print(helpText)
		default:
			fmt.Println("command", uinput, "not recognized")
		}

	}
	fmt.Println("BD:",bracketDepth)
	return false
}

func ifErrStackWarn(err error, neededVals int) bool { //check if there is an error, and if so, warn about the amout of values in stack
	if err != nil {
		fmt.Println("not enough values in stack, need at least", neededVals, "values")
		return true
	}

	return false
}

func pop2Vals(stk *stack.Stack) (val1 float64, val2 float64, err error) {
	val1, is1Float := stk.Pop().(float64)
	val2, is2Float := stk.Pop().(float64)
	if !is1Float || !is2Float {
		err = errors.New("pop values failed, not enough values or values are not numbers")
		return 0, 0, err
	}

	return val1, val2, nil
}

func showStack(stk *stack.Stack) {
	stackCopy := *stk
	for stackCopy.Peek() != nil {
		fmt.Println(stackCopy.Pop())
	}
}

func meditate() { //stops it from panicing
	if r := recover(); r != nil {
		fmt.Println("recovered from ", r)
	}
}

func getNumber(reader *bufio.Reader) (float64, error) {
	uinput, err := reader.ReadString('\n')
	ec(err)

	uinput = strings.Replace(uinput, "\n", "", -1)
	numinput, err := strconv.ParseFloat(uinput, 64)
	ec(err)

	return numinput, nil
}

func ec(err error) {
	if err != nil {
		panic(err)
	}
}
