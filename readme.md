# Invasion 
built with Go version 1.16.5

### Running the binary

run the command `make build` in project directory

run the binary with this command:
`./main -map YOUR_FILE.txt -n 5`

you can run the sample file with:
`./main -map map.txt -n 5`

#### Flags:
`-map (your_textFile)` - see below on how to form a well formed text input file

`-n (number)` - the number of aliens to spawn on the map. Note: at most can be 2*number of cities (in the maximum case, the game ends instantly)

the program will end by printing the cities and their paths left in the game, same format as input file.

### Running the tests
`make test`

### Assumptions
The project describes an example input file that has cities in its path not defined. I'm assuming that any city found
in a text file will have its own line describing its own paths. Each city will have at least 1 path. There are no sinks (in the beginning at least).

I assume the input file will be formed as such:
```
foo north=bar south=baz
bar south=foo
baz north=foo
```

### Notes
Initially, `MoveAlien` was `MoveAliens`, which would move all aliens at the same time.
Upon running the program, I noticed this could cause ping-pong behavior where two aliens could continuously swap positions.
Because of this, it only moves one alien at a time, and then checks for battle conditions

`SuperEpicAlienBattle` used to loop over all cities, however this felt like a waste of resources. Instead, since the above change
made aliens move one at a time, I decided to simply execute a check if they moved into a city with 1 alien, and execute
a fight after moving the 2nd alien in.
