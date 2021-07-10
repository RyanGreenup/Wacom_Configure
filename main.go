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

// Global Variables
const DEBUG bool = false
const DEVEL bool = false
const n int = 3 // How many buttons are there on the wacom? (The fourth button isn't detected by `xsetwacom set $id Button 4 ...`)

// NOTE in python
// import pyautogui
// pyautogui.position()

func main() {

	configure_stylus_area()
	configure_pad_buttons()

}

func configure_stylus_area() {
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
	cs := regexp.MustCompile(`\s`).Split(call_string, -1)

	if DEBUG {
		fmt.Println(call_string)
		// cs := strings.Split(call_string, " ")
		fmt.Println(cs)
	}

	// Finally use that string to adjust the Wacom Tablet
	_, err := exec.Command(cs[0], cs[1:]...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}

}

func configure_pad_buttons() {

	if DEVEL {
		fmt.Println("Use `xsetwacom --list` and record the $id of the PAD")
	}

	button_commands, err := get_wacom_button_exec_string()
	if err != nil {
		fmt.Println("Unable to generate commands")
	}
	// button_commands = [n]string{"xsetwacom set 16 Button 1 \"key +ctrl +shift p -ctrl -shift\"", "xsetwacom set 16 Button 1 \"key +ctrl +shift p -ctrl -shift\"", "xsetwacom set 16 Button 1 \"key +ctrl +shift p -ctrl -shift\""}

	if err != nil {
		fmt.Println("Unable to generate commands for buttons")
		fmt.Println(err)

	}

	if DEBUG {
		fmt.Println("Outside Loop")
		fmt.Println(button_commands)
	}

	for key, command_string := range button_commands {
		cs := regexp.MustCompile(`\s`).Split(command_string, -1)
		if DEBUG {
			fmt.Println("Inside Loop, iteration:", key)
			fmt.Println(cs)
		} else {
			_ = key
		}
		// Finally use that string to adjust the Wacom Tablet
		_, err := exec.Command(cs[0], cs[1:]...).Output()
		if err != nil {
			fmt.Printf("%s", err)
		}
	}

}

func get_wacom_button_exec_string() ([3]string, error) {

	// TODO should this be specified from the command line
	// No, you're just recreating the CLI, advice the user
	var keybindings [n]string = [n]string{"key +ctrl +shift p -ctrl -shift", "key +ctrl +shift = = = -ctrl -shift", "key +ctrl minus minus minus -ctrl "} // 3 is RMB

	// TODO should this be a variable?
	var commands [n]string = [n]string{"echo Foo", "echo blah", "echo meh"}

	for i := 0; i < n; i++ {

		var s string = ""
		// Try and get the ID or just print the devices
		id_val, err := get_wacom_id("PAD")
		if err != nil {
			fmt.Println("Unable to get ID value for PAD, try running: ")
			fmt.Println("xsetwacom --list")
			return [n]string{"", "", ""}, err
		} else {
			s = strings.Join(
				[]string{
					"xsetwacom set ",
					fmt.Sprint(id_val),
					" Button ",
					fmt.Sprint(i + 1),
					fmt.Sprint(" " + keybindings[i]),
				}, "")
			commands[i] = s
		}

	}

	print_button_commands(commands)

	return commands, nil
}

func print_button_commands(commands [n]string) {
	fmt.Println("The following commands were executed for the buttons" + "\n" +
		"the same commands can freely be re-executed to modify the mappings" + "\n")
	for key, command := range commands {
		_ = key
		fmt.Println(command)
	}
	fmt.Println("\n")
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
	id_val, err := get_wacom_id("STYLUS")
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
func get_wacom_id(device_name string) (int64, error) {

	out, err := exec.Command("xsetwacom", "--list").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	out_lines := strings.Split(string(out), "\n")
	for _, line := range out_lines {
		if strings.Contains(line, device_name) {

			// line = strings.Replace(line, ".*id", "foo", 1)

			re := regexp.MustCompile(`.*id:\ `)
			line := re.ReplaceAllString(line, "")

			re = regexp.MustCompile(`[\s]*type: ` + device_name)
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
