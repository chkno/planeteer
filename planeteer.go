/* Planeteer: Give trade route advice for Planets: The Exploration of Space
 * Copyright (C) 2011  Scott Worley <sworley@chkno.net>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import "flag"
import "fmt"
import "json"
import "os"
import "strings"

var start = flag.String("start", "",
	"The planet to start at")

var end = flag.String("end", "",
	"A comma-separated list of acceptable ending planets.")

var planet_data_file = flag.String("planet_data_file", "planet-data",
	"The file to read planet data from")

var fuel = flag.Int("fuel", 16, "Reactor units")

var hold = flag.Int("hold", 300, "Size of your cargo hold")

var start_edens = flag.Int("start_edens", 0,
	"How many Eden Warp Units are you starting with?")

var end_edens = flag.Int("end_edens", 0,
	"How many Eden Warp Units would you like to keep (not use)?")

var cloak = flag.Bool("cloak", false,
	"Make sure to end with a Device of Cloaking")

var drones = flag.Int("drones", 0, "Buy this many Fighter Drones")

var batteries = flag.Int("batteries", 0, "Buy this many Shield Batterys")

var visit_string = flag.String("visit", "",
	"A comma-separated list of planets to make sure to visit")

func visit() []string {
	return strings.Split(*visit_string, ",")
}

type Commodity struct {
	BasePrice int
	CanSell   bool
	Limit     int
}
type Planet struct {
	BeaconOn bool
	/* Use relative prices rather than absolute prices because you
	   can get relative prices without traveling to each planet. */
	RelativePrices map[string]int
}
type planet_data struct {
	Commodities map[string]Commodity
	Planets     map[string]Planet
	pi, ci      map[string]int // Generated; not read from file
}

func ReadData() (data planet_data) {
	f, err := os.Open(*planet_data_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		panic(err)
	}
	return
}

/* This program operates by filling in a state table representing the best
 * possible trips you could make; the ones that makes you the most money.
 * This is feasible because we don't look at all the possible trips.
 * We define a list of things that are germane to this game and then only
 * consider the best outcome in each possible game state.
 *
 * Each cell in the table represents a state in the game.  In each cell,
 * we track two things: 1. the most money you could possibly have while in
 * that state and 2. one possible way to get into that state with that
 * amount of money.
 *
 * A basic analysis can be done with a two-dimensional table: location and
 * fuel.  planeteer-1.0 used this two-dimensional table.  This version
 * adds features mostly by adding dimensions to this table.
 *
 * Note that the sizes of each dimension are data driven.  Many dimensions
 * collapse to one possible value (ie, disappear) if the corresponding
 * feature is not enabled.
 */

// The official list of dimensions:
const (
	// Name                Num  Size  Description
	Edens        = iota //   1     3  # of Eden warp units (0 - 2 typically)
	Cloaks              //   2     2  # of Devices of Cloaking (0 or 1)
	UnusedCargo         //   3     4  # of unused cargo spaces (0 - 3 typically)
	Fuel                //   4    17  Reactor power left (0 - 16)
	Location            //   5    26  Location (which planet)
	Hold                //   6    15  Cargo bay contents (a *Commodity or nil)
	NeedFighters        //   7     2  Errand: Buy fighter drones (needed or not)
	NeedShields         //   8     2  Errand: Buy shield batteries (needed or not)
	Visit               //   9  2**N  Visit: Stop by these N planets in the route

	NumDimensions
)

func bint(b bool) int {
	if b {
		return 1
	}
	return 0
}

func DimensionSizes(data planet_data) []int {
	eden_capacity := data.Commodities["Eden Warp Units"].Limit
	cloak_capacity := bint(*cloak)
	dims := make([]int, NumDimensions)
	dims[Edens] = eden_capacity + 1
	dims[Cloaks] = cloak_capacity + 1
	dims[UnusedCargo] = eden_capacity + cloak_capacity + 1
	dims[Fuel] = *fuel + 1
	dims[Location] = len(data.Planets)
	dims[Hold] = len(data.Commodities)
	dims[NeedFighters] = bint(*drones > 0) + 1
	dims[NeedShields] = bint(*batteries > 0) + 1
	dims[Visit] = 1 << uint(len(visit()))
	return dims
}

func StateTableSize(dims []int) int {
	sum := 0
	for _, size := range dims {
		sum += size
	}
	return sum
}

type State struct {
	funds, from int
}

func NewStateTable(dims []int) []State {
	return make([]State, StateTableSize(dims))
}

func EncodeIndex(dims, addr []int) int {
	index := addr[0]
	for i := 1; i < len(dims); i++ {
		index = index*dims[i] + addr[i]
	}
	return index
}

func DecodeIndex(dims []int, index int) []int {
	addr := make([]int, len(dims))
	for i := len(dims) - 1; i > 0; i-- {
		addr[i] = index % dims[i]
		index /= dims[i]
	}
	addr[0] = index
	return addr
}

/* What is the value of hauling 'commodity' from 'from' to 'to'?
 * Take into account the available funds and the available cargo space. */
func TradeValue(data planet_data,
from, to Planet,
commodity string,
initial_funds, max_quantity int) int {
	if !data.Commodities[commodity].CanSell {
		return 0
	}
	from_relative_price, from_available := from.RelativePrices[commodity]
	if !from_available {
		return 0
	}
	to_relative_price, to_available := to.RelativePrices[commodity]
	if !to_available {
		return 0
	}

	base_price := data.Commodities[commodity].BasePrice
	from_absolute_price := from_relative_price * base_price
	to_absolute_price := to_relative_price * base_price
	buy_price := from_absolute_price
	sell_price := int(float64(to_absolute_price) * 0.9)
	var can_afford int = initial_funds / buy_price
	quantity := can_afford
	if quantity > max_quantity {
		quantity = max_quantity
	}
	return (sell_price - buy_price) * max_quantity
}

func FindBestTrades(data planet_data) [][]string {
	// TODO: We can't cache this because this can change based on available funds.
	best := make([][]string, len(data.Planets))
	for from := range data.Planets {
		best[data.pi[from]] = make([]string, len(data.Planets))
		for to := range data.Planets {
			best_gain := 0
			price_list := data.Planets[from].RelativePrices
			if len(data.Planets[to].RelativePrices) < len(data.Planets[from].RelativePrices) {
				price_list = data.Planets[to].RelativePrices
			}
			for commodity := range price_list {
				gain := TradeValue(data,
					data.Planets[from],
					data.Planets[to],
					commodity,
					10000000,
					1)
				if gain > best_gain {
					best[data.pi[from]][data.pi[to]] = commodity
					gain = best_gain
				}
			}
		}
	}
	return best
}

// (Example of a use case for generics in Go)
func IndexPlanets(m *map[string]Planet) map[string]int {
	index := make(map[string]int, len(*m))
	i := 0
	for e := range *m {
		index[e] = i
		i++
	}
	return index
}
func IndexCommodities(m *map[string]Commodity) map[string]int {
	index := make(map[string]int, len(*m))
	i := 0
	for e := range *m {
		index[e] = i
		i++
	}
	return index
}

func main() {
	flag.Parse()
	data := ReadData()
	data.pi = IndexPlanets(&data.Planets)
	data.ci = IndexCommodities(&data.Commodities)
	dims := DimensionSizes(data)
	table := NewStateTable(dims)
	table[0] = State{1, 1}
	best_trades := FindBestTrades(data)

	for from := range data.Planets {
		for to := range data.Planets {
			best_trade := "(nothing)"
			if best_trades[data.pi[from]][data.pi[to]] != "" {
				best_trade = best_trades[data.pi[from]][data.pi[to]]
			}
			fmt.Printf("%s to %s: %s\n", from, to, best_trade)
		}
	}
}
