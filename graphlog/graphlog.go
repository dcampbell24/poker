package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Chip size of the small bet.
const _SB = 10

const wxtErr = `Error when returning from gnuplot:
If you didn't see dashed lines, you lack support for the wxt terminal.
If you saw dashed lines, something else went wrong while running gnuplot.`

func graphData() error {
	var count int
	sums  := make(map[string]int)
	index := make([]string, 0)

	in, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer in.Close()
	bufin := bufio.NewReader(in)

	tmp, err := ioutil.TempFile(".", "acpcLog")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	for {
		line, err := bufin.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		line = strings.TrimRight(line, "\r\n")

		if line[:len("STATE")] == "STATE" {
			state := strings.Split(line, ":")
			scores := strings.Split(state[4], "|")
			players := strings.Split(state[5], "|")
			for i, p := range players {
				score, err := strconv.Atoi(scores[i])
				if err != nil {
					return err
				}
				sums[p] += score
			}
			if count == 0 {
				for _, p := range players {
					index = append(index, p)
				}
				fmt.Fprintln(tmp, "# Players:", index)
			}
			count++
			fmt.Fprintf(tmp, "%d ", count)
			for _, p := range index {
				fmt.Fprintf(tmp, "%f ", float64(sums[p])/_SB)
			}
			fmt.Fprintln(tmp)
		}
	}
	gp := exec.Command("gnuplot")
	gpipe, err := gp.StdinPipe()
	if err != nil {
		return err
	}
	defer gpipe.Close()

	err = gp.Start()
	if err != nil {
		return err
	}

	// Do stuff here!
	fmt.Fprintf(gpipe, "set terminal wxt dashed;")
	fmt.Fprintf(gpipe, `set xlabel "Hand Number";`)
	fmt.Fprintf(gpipe, `set ylabel "Score [SB]";`)
	fmt.Fprintf(gpipe, `set title "%s";`, in.Name())
	cmd := "plot "
	for i, p := range index[:len(index)-1] {
		cmd += fmt.Sprintf(`"%s" u 1:%d title "%s" w l, `, tmp.Name(), i+2, p)
	}
	cmd += fmt.Sprintf(`"%s" u 1:%d title "%s" w l`, tmp.Name(), len(index)+1, index[len(index)-1])
	fmt.Fprintln(gpipe, cmd)

	stdin := bufio.NewReader(os.Stdin)
	for {
		c, err := stdin.ReadByte()
		if err != nil {
			return err
		}
		if c == '\n' {
		fmt.Fprintln(gpipe, "quit")
		break
		}
	}
	// End stuff

	err = gp.Wait()
	if err != nil {
		return fmt.Errorf(wxtErr)
	}
	return nil
}

func main() {
	err := graphData()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
