evolver
=======

Go implementation of evolver. 


How to run
===========
To build the Evolver simulator, you need SDL library (including development headers).



On Ubuntu:

    sudo apt-get install libsdl2-dev libsdl2-image-dev libsdl2-mixer-dev libsdl2-ttf-dev
    go get
    go build
    ./evolver

On Windows (not testet yet):
Point your browser to: [https://www.libsdl.org/download-2.0.php](https://www.libsdl.org/download-2.0.php) and download development libraries.

Open some civilized terminal and go to wherever you cloned the repository. (Make sure you have GOPATH set).

    go get
    go build
    ./evolver.exe
    
Usage
===========
For now, there are only 3 keys needed to operate the simulation:

	'+' make it faster
	'-' make it slower
	'esc' quit the simulation
