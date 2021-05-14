package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
)

// NOTE in python
// import pyautogui
// pyautogui.position()

func main() {

	main_query()

}

func main_query() {
	DEBUG := false
	DEVEL := false
	// 	robotgo.ScrollMouse(10, "up")
	// 	robotgo.MouseClick("left", true)
	// NOTE use `xdotool getmouselocation`
	// 	robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
	fmt.Println("Move the mouse to the bottom left of the writing area")
	press_enter()
	xbl, ybl := robotgo.GetMousePos()

	fmt.Println("Move the mouse to the top left of the writing area")
	press_enter()
	xtl, ytl := robotgo.GetMousePos()

	if DEVEL {
		// TODO this should use those functions to get the values
		fmt.Println("Use `xsetwacom --list` and record the $id of the STYLUS")
		fmt.Println("Get the Ratio of the Tablet with `xsetwacom get $id Area` (mine is 1.6)")
	}

	if DEBUG {
		fmt.Println(xbl, ybl)
		fmt.Println(xtl, ytl)
		height := ybl - ytl
		ratio := 1.6
		fmt.Printf("%dx%d+%d+%d\n", int(float64(height)*ratio), height, int(xbl), int(ytl))
	}

	call_string := get_wacom_exec_string(xtl, ytl, xbl, ybl)
	fmt.Println(call_string)

}

func get_wacom_exec_string(xtl int, ytl int, xbl int, ybl int) string {

	height := ybl - ytl
	ratio := get_wacom_ratio()

	if height < 0 {
		fmt.Println("Did you get that backwards?")
		// TODO use a recursive function to requery
	}

	// Try and get the ID or just print the devices
	var s string = ""
	id_val, err := get_wacom_id()
	if err != nil {
		print_wacom_devices()
		fmt.Printf("xsetwacom set $id MapToOutput ")
		s = strings.Join(
			[]string{
				"xsetwacom set $id MapToOutput",
			}, "")
	} else {
		s = strings.Join(
			[]string{
				"xsetwacom set ",
				fmt.Sprint(id_val),
				" MapToOutput ",
			}, "")
	}

	s = strings.Join(
		[]string{
			"xsetwacom set ",
			fmt.Sprint(id_val),
			" MapToOutput ",
		}, "")

	s = strings.Join(
		[]string{
			s,
			fmt.Sprint(int(float64(height) * float64(ratio))),
			"x",
			fmt.Sprint(height),
			"+",
			fmt.Sprint(int(xbl)),
			"+",
			fmt.Sprint(int(ytl)),
		}, "")

	return s
}

func press_enter() {
	buf := bufio.NewReader(os.Stdin)
	fmt.Println("Press Enter to Continue")
	fmt.Printf("> ")
	_, err := buf.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
	}
}

// This uses Regex to get the id value, it's not very
// robust and should atleast have error handling before it's properly
// implemented.
// Python might be better for this to be honest.
func get_wacom_id() (int64, error) {

	out, err := exec.Command("xsetwacom", "--list").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	out_lines := strings.Split(string(out), "\n")
	for _, line := range out_lines {
		if strings.Contains(line, "STYLUS") {

			// line = strings.Replace(line, ".*id", "foo", 1)

			re := regexp.MustCompile(`.*id:\ `)
			line := re.ReplaceAllString(line, "")

			re = regexp.MustCompile(`[\s]*type: STYLUS`)
			line = re.ReplaceAllString(line, "")

			re = regexp.MustCompile(`[\s]*`)
			line = re.ReplaceAllString(line, "")

			id_val, _ := strconv.ParseInt(line, 10, 32)
			return id_val, nil

			// re := regexp.MustCompile(`(?<=id:\ )\d`)
			// fmt.Printf("%q\n", re.Find([]byte(line)))
		}
	}

	return 0, errors.New("Could not extract id value")

}

func print_wacom_devices() {
	out, err := exec.Command("xsetwacom", "--list").Output()
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Println(string(out))
	}
	fmt.Println("Use the ID from the type: STYLUS")
}

func get_wacom_ratio() float64 {
	ratio := 1.6
	return ratio
}
