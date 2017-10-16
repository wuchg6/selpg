package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"bufio"
	"io"
	"io/ioutil"
)

var (
	maxInt int = 1 << 32 - 1
	startPage int
	endPage int
	inFilename string
	pageLen int = 72
	pageType int = 'l'
	printDest string
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"\nUSAGE: %v -sstartPage -eendPage [ -f | -llinesPerPage ] [ -ddest ] [ inFilename ]\n",
		os.Args[0])
}

func processArgs() {
	if (len(os.Args) < 3) {
		fmt.Fprintf(os.Stderr, "%v: not enough arguments\n", os.Args[0])
		usage()
		os.Exit(1)
	}

	/*
	 * handle 1st arg
	 * desired form: -sstartPage
	 */
	s1 := os.Args[1]
	if (s1[:2] != "-s") {
		fmt.Fprintf(os.Stderr, "%v: 1st arg should be -sstartPage\n", os.Args[0])
		usage()
		os.Exit(2)
	}
	startPage, _ = strconv.Atoi(s1[2:])
	if (startPage < 1 || startPage > maxInt - 1) {
		fmt.Fprintf(os.Stderr, "%v: invalid start page %v\n", os.Args[0], s1[2:])
		usage()
		os.Exit(3)
	}

	/*
	 * handle 2nd arg
	 * desired form: -eendPage
	 */
	s2 := os.Args[2]
	if (s2[:2] != "-e") {
		fmt.Fprintf(os.Stderr, "%v: 2st arg should be -eendPage\n", os.Args[0])
		usage()
		os.Exit(4)
	}
	endPage, _ = strconv.Atoi(s2[2:])
	if (endPage < 1 || endPage > maxInt - 1 || endPage < startPage) {
		fmt.Fprintf(os.Stderr, "%v: invalid end page %v\n", os.Args[0], s2[2:])
		usage()
		os.Exit(5)
	}

	/*
	 * handle optional args
	 * [ -f | -llinesPerPage ] [ -ddest ]
	 */
	argPos := 3
	for (argPos < len(os.Args) && os.Args[argPos][0] == '-') {
		s := os.Args[argPos]
		switch s[1] {
		case 'l':
			pageLen, _ = strconv.Atoi(s[2:])
			if (pageLen < 1 || pageLen > maxInt - 1) {
				fmt.Fprintf(os.Stderr, "%v: invalid page length %v\n", os.Args[0], s[2:])
				usage()
				os.Exit(6)
			}
			if (pageType == 'f') {
				fmt.Fprintf(os.Stderr, "%v: could not have both -f and -llinesPerPage", os.Args[0])
				usage()
				os.Exit(6)
			}
			argPos++

		case 'f':
			if (s != "-f") {
				fmt.Fprint(os.Stderr, "%v: option should be \"-f\"\n", os.Args[0])
				usage()
				os.Exit(7)
			}
			if (pageLen != 72) {
				fmt.Fprintf(os.Stderr, "%v: could not have both -f and -llinesPerPage", os.Args[0])
				usage()
				os.Exit(7)
			}
			pageType = 'f'
			argPos++

		case 'd':
			if (s == "-d") {
				fmt.Fprintf(os.Stderr, "%v: -d option requires a printer destination\n", os.Args[0])
				usage()
				os.Exit(8)
			}
			printDest = s[2:]
			argPos++
		
		default:
			fmt.Fprintf(os.Stderr, "%v: unknown option %v\n", os.Args[0], s)
			usage();
			os.Exit(9)
		}
	}

	/*
	 * handle optional args
	 * [ inFilename ]
	 */
	if (argPos < len(os.Args)) {
		inFilename = os.Args[argPos]
		_, err := os.Stat(inFilename)
		if (os.IsNotExist(err)) {
			fmt.Fprintf(os.Stderr, "%v: input file \"%v\" does not exist\n", os.Args[0], inFilename)
			os.Exit(10)
		}
	}
}

func processInput() {
	/*
	 * set the input source
	 */
	fin := os.Stdin
	if (len(inFilename) != 0) {
		f, e := os.Open(inFilename)
		if (e != nil) {
			fmt.Fprintf(os.Stderr, "%v: could not open input file \"%v\"\n", os.Args[0], inFilename)
			os.Exit(11)
		}
		fin = f
	}
	reader := bufio.NewReader(fin)

	/*
	 * set the output destination
	 */
	fout := os.Stdout
	var cmd *exec.Cmd = nil
	var cmdIn io.WriteCloser = nil
	var cmdOut io.ReadCloser = nil
	var cmdErr io.ReadCloser = nil
	if (len(printDest) != 0) {
		cmd = exec.Command("lp", "-d", printDest)
		dataIn, e1 := cmd.StdinPipe()
		dataOut, e2 := cmd.StdoutPipe()
		dataErr, e3 := cmd.StderrPipe()
		if (e1 != nil && e2 != nil && e3 != nil) {
			fmt.Fprintf(os.Stderr, "%v: could not open pipe to \"%v\"\n", os.Args[0], printDest)
			os.Exit(12)
		}
		cmdIn = dataIn
		cmdOut = dataOut
		cmdErr = dataErr
	}
	writer := bufio.NewWriter(fout)

	/*
	 * start the desired command
	 * then, we can input data through pipe
	 */
	if (cmd != nil) {
		cmd.Start()
	}

	if (pageType == 'l') {
		lineCount := 0
		pageCount := 1

		for  {
			line, _, _ := reader.ReadLine()

			if (line == nil) {
				break
			}

			lineCount++
			if (lineCount > pageLen) {
				pageCount++
				lineCount = 0
			}

			if (pageCount >= startPage && pageCount <= endPage) {
				if (cmdIn == nil) {
					writer.Write(line)
					writer.WriteByte('\n')
				} else {
					cmdIn.Write(line)
					cmdIn.Write([]byte{'\n'})
				}
			}
		}
	} else {
		pageCount := 1

		for {
			c, err := reader.ReadByte()

			if (err == io.EOF) {
				break
			}

			if (c == '\f') {
				pageCount++
			}

			if (pageCount >= startPage && pageCount <= endPage) {
				writer.WriteByte(c)
				if (cmdIn == nil) {
					writer.WriteByte(c)
				} else {
					cmdIn.Write([]byte{c})
				}
			}
		}
	}

	writer.Flush()

	/*
	 * process the information of lp -d 
	 */
	if (cmd != nil) {
		cmdIn.Close()

		var lpOut []byte

		lpOut, _ = ioutil.ReadAll(cmdOut)
		cmdOut.Close()
		if (len(lpOut) != 0) {
			fmt.Println(string(lpOut))
		}

		lpOut, _ = ioutil.ReadAll(cmdErr)
		cmdErr.Close()
		if (len(lpOut) != 0) {
			fmt.Println(string(lpOut))
		}

		cmd.Wait()
	}

	fin.Close()
	fout.Close()
}

func main() {
	processArgs()
	processInput()
}
