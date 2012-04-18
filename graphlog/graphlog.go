package main

import (
	"bufio"
	"fmt"
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

	log, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return err
	}

	tmp, err := ioutil.TempFile(".", "acpcLog")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	bufout := bufio.NewWriter(tmp)

	lines := strings.Split(string(log), "\n")
	for _, line := range lines {
		if len(line) > 4 && line[:5] == "STATE" {
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
				fmt.Fprintln(bufout, "# Players:", index)
			}
			count++
			fmt.Fprintf(bufout, "%d ", count)
			for _, p := range index {
				fmt.Fprintf(bufout, "%f ", float64(sums[p])/_SB)
			}
			fmt.Fprintln(bufout)
		}
	}
	bufout.Flush()
	tmp.Close()
	// Plotting Code.
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
	cmd := "set t wxt dash;"
	cmd += `set xl "Hand Number";`
	cmd += `set yl "Score [SB]";`
	cmd += fmt.Sprintf(`set title "%s";`, os.Args[1])
	cmd += "p "
	for i, p := range index[:len(index)-1] {
		cmd += fmt.Sprintf(`"%s" u 1:%d title "%s" w l, `, tmp.Name(), i+2, p)
	}
	cmd += fmt.Sprintf(`"" u 1:%d title "%s" w l`, len(index)+1, index[len(index)-1])
	fmt.Fprintln(gpipe, cmd)
	fmt.Println("Press _Enter_ to quit.")

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
