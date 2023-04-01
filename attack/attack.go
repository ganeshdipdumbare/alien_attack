package attack

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
)

// Alien struct to hold alien details
// name: name of the alien
// cityName: name of the city the alien is in
// noOfVisits: number of times the alien has visited a city
// isDead: flag to check if the alien is dead
// isStuck: flag to check if the alien is stuck
type Alien struct {
	name       string
	cityName   string
	noOfVisits int
	isDead     bool
	isStuck    bool
}

// City struct to hold city details
// name: name of the city
// connectedCities: list of connected cities
// aliens: list of aliens in the city
// isDestroyed: flag to check if the city is destroyed
type City struct {
	name            string
	connectedCities []string
	aliens          []*Alien
	isDestroyed     bool
}

// World struct to hold world details
// cities: map of cities in the world
// lock: lock to make the world thread safe
type World struct {
	cities map[string]*City
	lock   *sync.RWMutex
}

// CreateWorld function to create world from the map file
// mapFileName: name of the file containing the map of the world
// returns: world object and error if any
func CreateWorld(mapFileName string) (*World, error) {
	world := &World{
		cities: make(map[string]*City, 0),
		lock:   &sync.RWMutex{},
	}

	f, err := os.Open(mapFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, " ")
		if len(data) > 0 {
			cityName := data[0]
			connectedCities := data[1:]
			world.cities[cityName] = &City{
				name:            cityName,
				connectedCities: connectedCities,
				aliens:          make([]*Alien, 0),
			}
		}
	}

	// add connected cities in the world if they are not already there
	for _, city := range world.cities {
		for _, connectedCity := range city.connectedCities {
			connectedCityName := strings.Split(connectedCity, "=")[1]
			if _, ok := world.cities[connectedCityName]; !ok {
				world.cities[connectedCityName] = &City{
					name:            connectedCity,
					connectedCities: make([]string, 0),
					aliens:          make([]*Alien, 0),
				}
			}
		}
	}

	return world, nil
}

// PrintMap function to print the map of the world
func (w *World) PrintMap() {
	// write cities and connected cities to file per line
	worldMap := "\n"
	for _, city := range w.cities {
		if !city.isDestroyed {
			worldMap = fmt.Sprintf("%v %v", worldMap, city.name)
			for _, connectedCity := range city.connectedCities {
				worldMap = fmt.Sprintf("%v %v", worldMap, connectedCity)
			}
			worldMap = worldMap + "\n"
		}
	}
	fmt.Println(worldMap)
}

// GetNoOfCities function to get number of cities from the world
func (w *World) GetNoOfCities() int {
	return len(w.cities)
}

// getNoOfDestroyedCities get no of destroyed cities
func (w *World) getNoOfDestroyedCities() int {
	noOfDestroyedCities := 0
	for _, city := range w.cities {
		if city.isDestroyed {
			noOfDestroyedCities++
		}
	}
	return noOfDestroyedCities
}

// create aliens with random cities
// no of aliens should be less than or equal to no of cities
func createAliensWithRandomCity(w *World, noOfAliens int) []*Alien {
	aliens := make([]*Alien, 0)
	for i := 0; i < noOfAliens; i++ {
		aliens = append(aliens, &Alien{
			name:       fmt.Sprintf("Alien-%v", i),
			noOfVisits: 0,
			isDead:     false,
			isStuck:    false,
		})
	}

	i := noOfAliens
	for k := range w.cities {
		if i == 0 {
			break
		}

		addAlienToCity(w.cities[k], aliens[i-1])
		i--
	}

	return aliens
}

func (w *World) UnleashAliens(noOfAliens int) {
	aliens := createAliensWithRandomCity(w, noOfAliens)
	var wg sync.WaitGroup
	for {

		// no lock is required here as the aliens are not modifying the world
		// at this point
		isOver, liveAliens := isAttackOver(w, len(aliens))
		if isOver {
			break
		}

		for _, a := range liveAliens {
			wg.Add(1)
			go func(alien *Alien) {
				alien.VisitRandomConnectedCity(w)
				wg.Done()
			}(a)
		}

		wg.Wait()
	}
}

// func to add aliens to cities
func addAlienToCity(c *City, a *Alien) {
	a.cityName = c.name
	c.aliens = append(c.aliens, a)
}

// VisitRandomConnectedCity function to visit a random connected city
// w: world object
// a: alien object
func (a *Alien) VisitRandomConnectedCity(w *World) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// current city will be the previous city as the alien is moving to a new city
	previousCity := w.cities[a.cityName].name

	if a.isDead || a.isStuck {
		return
	}

	if len(w.cities[a.cityName].connectedCities) == 0 {
		a.isStuck = true
		return
	}

	connectedCitiesForAlien := w.cities[a.cityName].connectedCities
	randomCity := connectedCitiesForAlien[rand.Intn(len(connectedCitiesForAlien))]
	randomCityName := strings.Split(randomCity, "=")[1]
	a.cityName = randomCityName
	a.noOfVisits++

	w.cities[randomCityName].aliens = append(w.cities[randomCityName].aliens, a)
	if len(w.cities[randomCityName].aliens) > 1 {
		// kill aliens and destroy city as there are more than 1 alien in the city
		for _, v := range w.cities[randomCityName].aliens {
			v.isDead = true
		}

		w.cities[randomCityName].isDestroyed = true
		fmt.Printf("%v has been destroyed by %v and %v!\n", randomCityName, w.cities[randomCityName].aliens[0].name, w.cities[randomCityName].aliens[1].name)

		// remove the city from connected cities of other cities
		for _, city := range w.cities {
			for i, connectedCity := range city.connectedCities {
				if strings.Split(connectedCity, "=")[1] == randomCityName {
					city.connectedCities = append(city.connectedCities[:i], city.connectedCities[i+1:]...)
				}
			}
		}
	}

	if len(w.cities[randomCityName].connectedCities) == 0 {
		a.isStuck = true
	}

	// remove the alien from the previous city
	for i, alien := range w.cities[previousCity].aliens {
		if alien.name == a.name {
			w.cities[previousCity].aliens = append(w.cities[previousCity].aliens[:i], w.cities[previousCity].aliens[i+1:]...)
		}
	}
}

// isAttackOver function to check if the attack is over
// attack is over if all aliens are dead or stuck or visited 10000 cities
func isAttackOver(w *World, totalAliens int) (bool, []*Alien) {
	noOfStoppableAliens := 0
	liveAliens := make([]*Alien, 0)

	for _, city := range w.cities {
		for _, alien := range city.aliens {
			if alien.isDead || alien.isStuck {
				noOfStoppableAliens++
				continue
			}
			if alien.noOfVisits >= 10000 {
				noOfStoppableAliens++
				continue
			}
			liveAliens = append(liveAliens, alien)
		}
	}

	return noOfStoppableAliens == totalAliens, liveAliens
}
