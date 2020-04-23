# scalc
programable stack-based calculator written in go

the following commands are currently supported:

    q    quits the program
    l    view stack
    +    add last 2 numbers in the stack
    -    subtract the last number from the second-to-last number
    *    multiply the last 2 numbers in the stack
    /    divide the last number from the second-to-last number
    |    swap the last 2 values in stack
    $    store the second-to-last value under the index of the last value
    =    retrives the value stored under the index of the last value in stack
    ~    removes the last value
    _    duplicates the last value
    ?    if the last value is 0, skips to next line of instructions
    [    store instructions in stack until matching ] is reached
    #    if the last value was a set of instuctions stored with [, evaluate them
    h    print this helptext

enter any number to push it to the stack

multiple commands can be entered on a single line by seperating them with a space

any command line arguments will be parsed as if they were entered interactively.
