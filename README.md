# scalc
programable stack-based calculator written in go

the following commands are currently supported:

    q    quits the program
    =    view stack
    +    add last 2 numbers in the stack
    -    subtract the last number from the second-to-last number
    *    multiply the last 2 numbers in the stack
    /    divide the last number from the second-to-last number
    h    print this helptext

enter any number to push it to the stack

multiple commands can be entered on a single line by seperating them with a space

any command line arguments will be parsed as if they were entered interactively.
