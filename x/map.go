package x

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// MaxMoves is the maximum amount of moves an alien can make before being removed from the map
const MaxMoves = 10_000

// the allowed directions when building a map
var allowedDirections = map[string]bool{"north": true, "east": true, "south": true, "west": true}

type (
	// City has a name, paths to other cities, and alien invaders!
	City struct {
		name     string
		path     map[string]string // north=baz south=foo etc...
		invaders map[string]*Alien
	}

	// Map has many cities, and keeps track of all the invaders in the cities
	Map struct {
		cities map[string]*City  // baz:City
		aliens map[string]*Alien // alien-1:Alien
	}
)

// String prints the city as it would look in the input file
func (c City) String() string {
	str := fmt.Sprintf("%s", c.name)
	for direction, path := range c.path {
		str += fmt.Sprintf(" %s=%s", direction, path)
	}
	return str
}

// PrintCities prints the cities to the cli in the same format as the input file
func (m *Map) PrintCities() {
	for _, c := range m.cities {
		fmt.Println(c)
	}
}

// LoadMapFromSlice allows you to pass in a slice of strings to load the map, rather than a file
func LoadMapFromSlice(s []string) (*Map, error) {
	m := Map{
		cities: make(map[string]*City),
		aliens: make(map[string]*Alien),
	}
	for _, line := range s {
		lineSlice := strings.Split(line, " ")
		if len(lineSlice) < 2 {
			return nil, errors.New("not enough data in file: please define a filename and paths")
		}
		city, err := extractCity(lineSlice)
		if err != nil {
			return nil, err
		}
		m.cities[city.name] = city
	}

	return &m, nil
}

// LoadMapFromFile loads a map from a given file object
func LoadMapFromFile(f *os.File) (*Map, error) {
	m := Map{
		cities: make(map[string]*City),
		aliens: make(map[string]*Alien),
	}

	// scan in the files
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lineSlice := strings.Split(line, " ")
		if len(lineSlice) < 2 {
			return nil, errors.New("not enough data in file: please define a city name and at least one path")
		}

		city, err := extractCity(lineSlice)
		if err != nil {
			return nil, err
		}

		m.cities[city.name] = city
	}

	return &m, nil
}

// extracts a city given a string with the format north=baz (more generally - direction=cityname)
func extractCity(line []string) (*City, error) {
	city := City{}
	city.name = line[0]
	city.invaders = make(map[string]*Alien)
	city.path = make(map[string]string)

	paths := line[1:]
	for _, v := range paths {
		path := strings.Split(v, "=")
		if _, ok := allowedDirections[path[0]]; !ok {
			return nil, errors.New("invalid direction: must use either north, west, south, or east")
		}
		city.path[path[0]] = path[1]
	}
	return &city, nil
}

// SeedAliens randomly places aliens in the cities
func (m *Map) SeedAliens(numAliens int) error {
	if numAliens >= len(m.cities)*2+1 {
		return errors.New("cannot have more than 2*number of cities aliens")
	}
	for numAliens > 0 {
		for name, city := range m.cities {
			if len(city.invaders) < 2 {
				newAlien := Alien{
					name:     fmt.Sprintf("alien-%d", numAliens),
					location: name,
				}
				city.invaders[newAlien.name] = &newAlien
				m.aliens[newAlien.name] = &newAlien
				numAliens--
				fmt.Printf("%s has invaded %s!\n", newAlien.name, city.name)
				if numAliens <= 0 {
					break
				}
			}
		}
	}
	return nil
}

// MoveAlien moves an alien randomly along a path in it's current city. returns true if aliens moved, false otherwise.
//
// Initially this moved ALL aliens at the same time, however this can causes ping-pong errors where two aliens
// end up just playing musical chairs between cities. it now only moves one alien at a time.
func (m *Map) MoveAlien() (anyMoved bool) {
	for _, alien := range m.aliens {
		if alien.moves >= MaxMoves {
			// we now retire this alien from the game.
			ct := m.cities[alien.location]
			delete(ct.invaders, alien.name)
			delete(m.aliens, alien.name)
			continue
		}

		city := m.cities[alien.location]
		// lets cleanup the dead ends
		if len(city.path) == 0 {
			m.cleanUp(city)
		}

		for _, nextCity := range city.path {
			if len(m.cities[nextCity].invaders) < 2 {
				next := m.cities[nextCity] // get the next city

				alien.location = nextCity         // change the aliens location
				delete(city.invaders, alien.name) // remove the alien from the previous city's invader slice
				next.invaders[alien.name] = alien // update the new cities occupancy
				alien.moves++                     // update the aliens moves
				anyMoved = true                   // set the anyMoved flag
				if len(next.invaders) >= 2 {
					m.superEpicAlienBattle(next.name)
				}
				return true
			}
		}
	}
	return anyMoved
}

// CheckBattleConditionsAndExec checks the conditions of the map and executes battles.
// This is mostly used so you can still do weird stuff like spawn 2*city aliens and have the game implode instantly.
func (m *Map) CheckBattleConditionsAndExec() {
	for _, city := range m.cities {
		if len(city.invaders) < 2 {
			continue
		}
		m.superEpicAlienBattle(city.name)
	}
}

// superEpicAlienBattle executes a battle between two Aliens in a given city. A super epic alien battle ensues.
func (m *Map) superEpicAlienBattle(city string) {
	dc := m.cities[city]
	if len(dc.invaders) == 2 {
		m.cleanUp(m.cities[city])
	}
}

func (m *Map) cleanUp(city *City) {
	// cleanup a city that has either just had a super epic fight, or the city is isolated and can be removed.
	for _, v := range city.invaders {
		fmt.Printf("%s ", v.name)
		delete(m.aliens, v.name)
	}
	fmt.Printf("are dead. %s is destroyed!\n", city.name)
	cityName := city.name
	delete(m.cities, cityName)
	m.removeAllPaths(cityName)
}

// helper function to remove a city from a city's path map
func (m *Map) removeAllPaths(cityName string) {
	for _, city := range m.cities {
		for dir, path := range city.path {
			if path == cityName {
				delete(city.path, dir)
			}
		}
	}
}
