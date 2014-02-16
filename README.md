This project implements the common game sokoban with Google Go.
There is no GUI, but only text output and input.

What is needed?
===============
You need to install Google Go. It is usually in your distros repository, e.g. for Ubuntu: golang, for Arch: go


How can you run it?
===================
See the short introduction to setup your workspace:
http://golang.org/doc/code.html
Then, simply use go install github.com/g3force/Go_Sokoban

How to use?
===========
    ~> Go_Sokoban [-r] [-m] [-i] [-s] [-l <levelfile>] [-f <outputFrequency>] [-d <debuglevel>] [-p]
    -r to directly run the algorithm
    -m for finding more than one solution
    -i for information
    -s for straightAhead
    -l for levelfile
    -f for outputFrequency
    -d for debuglevel
    -p for printing Surface regularly
    the order of parameters does not matter
