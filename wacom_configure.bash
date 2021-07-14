# Shell Settings
set -o errexit # Abort on nonzero
set -o nounset # Abort on unbound var
set -o pipefail # Don't hide errors within pipes

RATIO="16/10" # Ratio of the Wacom Tablet

main() {
    # Get the ID Values ..........................................
    id_val_stylus="$(xsetwacom --list | grep STYLUS | sed -E  's/.*id:\ ([[:digit:]]+).*/\1/g')"
    id_val_pad="$(xsetwacom --list | grep PAD | sed -E  's/.*id:\ ([[:digit:]]+).*/\1/g')"

    # . Retreive Co-ordinates (from mouse) ......................
    echo "Move mouse to bottom right"
    press_key_to_continue
    xbr="$(get_x)"
    ybr="$(get_y)"
    
    echo -e "\nMove mouse to top right"
    press_key_to_continue
    xtr="$(get_x)"
    ytr="$(get_y)"

    # Calculate the xbl and ytl Coordinates and width 
    let height=${ytr}+${ybr}
    let width=${height}\*${RATIO}
    let xbl=${xbr}-${width}
    ytl="${ytr}"

    # Choose Button bindings
    button_1="key +ctrl +shift p -ctrl -shift" # use "3" form RMB
    button_2="key +ctrl +shift = = = -ctrl -shift"
    button_3="key +ctrl minus minus minus -ctrl"
    
    # Make the Commands 
    area_command="xsetwacom set ${id_val_stylus} MapToOutput ${width}x${height}+${xbl}+${ytl}"
    button_command_1="xsetwacom set ${id_val_pad} Button 1 \"${button_1}\""
    button_command_2="xsetwacom set ${id_val_pad} Button 1 \"${button_2}\""
    button_command_3="xsetwacom set ${id_val_pad} Button 1 \"${button_3}\""

    # Print the Commands
    echo -e "\n\n Run these commands to configure Wacom tablet \n\n"
    echo ${area_command}
    echo ${button_command_1}
    echo ${button_command_2}
    echo ${button_command_3}

   
}


# Helper Funtions .............................................
function press_key_to_continue() {
    read -n 1 -p "Press any key to Confirm"
}

function get_x() {
    xdotool getmouselocation | cut -f 1 -d ' ' | grep -P -o "(?<=x:)[\d]+"
}

function get_y() {
    xdotool getmouselocation | cut -f 2 -d ' ' | grep -Po '(?<=y:)[\d]+'
}


main "${@:-}"
