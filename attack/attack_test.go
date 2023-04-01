package attack

import (
	"reflect"
	"sync"
	"testing"
)

func TestCreateWorld(t *testing.T) {
	type args struct {
		mapFileName string
	}
	tests := []struct {
		name    string
		args    args
		want    *World
		wantErr bool
	}{
		{
			name: "should create world from map file",
			args: args{
				mapFileName: "./testutils/cities_test1.txt",
			},
			want: &World{
				cities: map[string]*City{
					"city1": {
						name:            "city1",
						connectedCities: []string{"east=city2", "west=city3"},
						aliens:          []*Alien{},
						isDestroyed:     false,
					},
					"city2": {
						name:            "city2",
						connectedCities: []string{"west=city1", "south=city3"},
						aliens:          []*Alien{},
						isDestroyed:     false,
					},
					"city3": {
						name:            "city3",
						connectedCities: []string{},
						aliens:          []*Alien{},
						isDestroyed:     false,
					},
				},
				lock: &sync.RWMutex{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateWorld(tt.args.mapFileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWorld() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateWorld() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorld_GetNoOfCities(t *testing.T) {
	type fields struct {
		cities map[string]*City
		lock   *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "should return number of cities in the world",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						name:            "city1",
						connectedCities: []string{"east=city2", "west=city3"},
						aliens:          []*Alien{},
						isDestroyed:     false,
					},
					"city2": {
						name:            "city2",
						connectedCities: []string{"west=city1", "south=city3"},
						aliens:          []*Alien{},
						isDestroyed:     false,
					},
					"city3": {
						name:            "city3",
						connectedCities: []string{},
						aliens:          []*Alien{},
						isDestroyed:     false,
					},
				},
				lock: &sync.RWMutex{},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &World{
				cities: tt.fields.cities,
				lock:   tt.fields.lock,
			}
			if got := w.GetNoOfCities(); got != tt.want {
				t.Errorf("World.GetNoOfCities() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorld_UnleashAliens(t *testing.T) {
	type fields struct {
		cities map[string]*City
		lock   *sync.RWMutex
	}
	type args struct {
		noOfAliens int
	}
	tests := []struct {
		name                        string
		fields                      fields
		args                        args
		expectedNoOfDestroyedCities int
	}{
		{
			name: "unleash only 1 alien in the world",
			fields: fields{
				cities: map[string]*City{
					"city1": {
						name:            "city1",
						connectedCities: []string{"east=city2", "west=city3"},
						aliens:          []*Alien{},
						isDestroyed:     false,
					},
					"city2": {
						name:            "city2",
						connectedCities: []string{"west=city1", "south=city3"},
						aliens:          []*Alien{},
						isDestroyed:     false,
					},
					"city3": {
						name:            "city3",
						connectedCities: []string{},
						aliens:          []*Alien{},
						isDestroyed:     false,
					},
				},
				lock: &sync.RWMutex{},
			},
			args: args{
				noOfAliens: 1,
			},
			expectedNoOfDestroyedCities: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &World{
				cities: tt.fields.cities,
				lock:   tt.fields.lock,
			}
			w.UnleashAliens(tt.args.noOfAliens)
			if got := w.getNoOfDestroyedCities(); got != tt.expectedNoOfDestroyedCities {
				t.Errorf("World.UnleashAliens() = %v, want %v", got, tt.expectedNoOfDestroyedCities)
			}
		})
	}
}
