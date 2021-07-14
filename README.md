
# Wacom Configuration for Linux

A program to help write the command for `xsetwacom` in order to configure a wacom tablet on linux (see [xsetwacom docs](https://github.com/linuxwacom/xf86-input-wacom/wiki/Dual-and-Multi-Monitor-Set-Up#maptooutput) and the [arch wiki article](https://wiki.archlinux.org/title/Wacom_tablet#Adjusting_aspect_ratios)) for more information.

## Usage

There are two working scripts [One in go](./main.go) and another in [bash](./wacom_configure.bash).

### Bash

This is more to illustrate how to use xsetwacom:

  1. Plugin wacom
  2. Confirm that it is detected by running `xsetwacom --list`
  3. Run one of the scripts with `bash ./wacom_configure.bash` or `go run ./main.go`
    1. the `bash` script will only print what the commands should be where as the *Go* script will run them.
  5. Go will run the appropriate commands 

## NB

The ratio of the wacom tablet is hardcoded as a global variable in the scripts, It will be necessary to change that accordingly (or circles will become ovals).

