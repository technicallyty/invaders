package x

import (
	"os"
	"testing"
)

func TestFight(t *testing.T) {
	m, err := LoadMapFromSlice([]string{"foo north=bar", "bar south=foo"})
	if err != nil {
		t.Fatal(err)
	}
	err = m.SeedAliens(2)
	if err != nil {
		t.Fatal("SeedAliens should not fail")
	}

	if len(m.aliens) != 2 {
		t.Fatal("aliens not seeded correctly")
	}

	// move will execute a battle as there is only one other place an alien could move, with 1 alien already there.
	m.MoveAlien()

	if len(m.aliens) != 0 {
		t.Fatal("all aliens should be dead")
	}
}

func TestNumAliens(t *testing.T) {
	m, err := LoadMapFromSlice([]string{"foo north=bar south=baz", "bar south=foo", "baz north=foo"})
	if err != nil {
		t.Fatal("valid input to loadmap should not fail")
	}
	numCities := len(m.cities)
	maxAliens := numCities * 2
	aboveMax := maxAliens + 1

	err = m.SeedAliens(aboveMax)
	if err == nil {
		t.Fatal("invalid amount given to seed")
	}

	err = m.SeedAliens(maxAliens)
	if err != nil {
		t.Fatal("amount should be allowed")
	}
}

func TestInput(t *testing.T) {
	testCases := []struct {
		input  []string
		expErr bool
	}{
		{
			input:  []string{"foo north=bar", "bar south=foo"},
			expErr: false,
		},
		{
			input:  []string{"foo bogus=blah", "fjw jfoief=fjief"},
			expErr: true,
		},
	}
	for _, tc := range testCases {
		_, err := LoadMapFromSlice(tc.input)
		if tc.expErr {
			if err == nil {
				t.Fatal("invalid input was accept in load function")
			}
		} else {
			if err != nil {
				t.Fatal("valid input returned an error")
			}
		}
	}
}

func TestExitsOnMaxMoves(t *testing.T) {
	m, err := LoadMapFromSlice([]string{"foo north=bar", "bar south=foo"})
	if err != nil {
		t.Fatal("valid input returned an error")
	}
	err = m.SeedAliens(1)
	if err != nil {
		t.Fatal("SeedAliens should not error")
	}

	for i := -1; i < MaxMoves; i++ {
		m.MoveAlien()
	}

	if len(m.aliens) != 0 {
		t.Fatal("alien did not retire after maximum allowed moves.")
	}
}

func TestCleanUp(t *testing.T) {
	m := Map{
		cities: make(map[string]*City),
		aliens: make(map[string]*Alien),
	}

	alien := Alien{
		name:     "alien-test1",
		location: "burp",
		moves:    0,
	}

	testCity := City{
		name:     "burp",
		path:     make(map[string]string),
		invaders: map[string]*Alien{alien.name: &alien},
	}
	m.cities[testCity.name] = &testCity
	m.aliens[alien.name] = &alien

	m.cleanUp(&testCity)

	if len(m.aliens) != 0 {
		t.Fatal("there should be no aliens left")
	}
	if len(m.cities) != 0 {
		t.Fatal("there should be no cities left")
	}
}

func TestPathRemoval(t *testing.T) {
	testData := []string{
		"foo north=bar east=baz south=bitcoin west=ether",
		"bar south=foo",
		"baz west=foo",
		"bitcoin north=foo",
		"ether east=foo",
	}

	m, err := LoadMapFromSlice(testData)
	if err != nil {
		t.Fatal("valid input returned an error")
	}

	m.removeAllPaths("foo")
	for _, city := range m.cities {
		for _, path := range city.path {
			if path == "foo" {
				t.Fatal("foo should've been removed from path")
			}
		}
	}
}

func TestReadFromFile(t *testing.T) {
	file, err := os.Open("map_test.txt")
	if err != nil {
		t.Fatal("error opening file")
	}
	m, err := LoadMapFromFile(file)
	if err != nil {
		t.Fatal("failed to open valid file")
	}

	if len(m.cities) != 5 {
		t.Fatal("5 cities present in file but not in map")
	}
}
