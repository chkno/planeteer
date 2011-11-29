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
import "runtime/pprof"
import "strings"

var funds = flag.Int("funds", 0,
	"Starting funds")

var start = flag.String("start", "",
	"The planet to start at")

var flight_plan_string = flag.String("flight_plan", "",
	"Your hyper-holes for the day, comma-separated.")

var end_string = flag.String("end", "",
	"A comma-separated list of acceptable ending planets.")

var planet_data_file = flag.String("planet_data_file", "planet-data",
	"The file to read planet data from")

var fuel = flag.Int("fuel", 16, "Hyper Jump power left")

var hold = flag.Int("hold", 300, "Size of your cargo hold")

var start_edens = flag.Int("start_edens", 0,
	"How many Eden Warp Units are you starting with?")

var end_edens = flag.Int("end_edens", 0,
	"How many Eden Warp Units would you like to keep (not use)?")

var cloak = flag.Bool("cloak", false,
	"Make sure to end with a Device of Cloaking")

var drones = flag.Int("drones", 0, "Buy this many Fighter Drones")

var batteries = flag.Int("batteries", 0, "Buy this many Shield Batterys")

var drone_price = flag.Int("drone_price", 0, "Today's Fighter Drone price")

var battery_price = flag.Int("battery_price", 0, "Today's Shield Battery price")

var visit_string = flag.String("visit", "",
	"A comma-separated list of planets to make sure to visit")

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var visit_cache []string

func visit() []string {
	if visit_cache == nil {
		if *visit_string == "" {
			return nil
		}
		visit_cache = strings.Split(*visit_string, ",")
	}
	return visit_cache
}

var flight_plan_cache []string

func flight_plan() []string {
	if flight_plan_cache == nil {
		if *flight_plan_string == "" {
			return nil
		}
		flight_plan_cache = strings.Split(*flight_plan_string, ",")
	}
	return flight_plan_cache
}

var end_cache map[string]bool

func end() map[string]bool {
	if end_cache == nil {
		if *end_string == "" {
			return nil
		}
		m := make(map[string]bool)
		for _, p := range strings.Split(*end_string, ",") {
			m[p] = true
		}
		end_cache = m
	}
	return end_cache
}

type Commodity struct {
	BasePrice int
	CanSell   bool
	Limit     int
}
type Planet struct {
	BeaconOn bool
	Private  bool
	/* Use relative prices rather than absolute prices because you
	   can get relative prices without traveling to each planet. */
	RelativePrices map[string]int
}
type planet_data struct {
	Commodities map[string]Commodity
	Planets     map[string]Planet
	p2i, c2i    map[string]int // Generated; not read from file
	i2p, i2c    []string       // Generated; not read from file
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
 *
 * The order of the dimensions in the list of constants below determines
 * their layout in RAM.  The cargo-based 'dimensions' are not completely
 * independent -- some combinations are illegal and not used.  They are
 * handled as three dimensions rather than one for simplicity.  Placing
 * these dimensions first causes the unused cells in the table to be
 * grouped together in large blocks.  This keeps them from polluting
 * cache lines, and if they are large enough, allows the memory manager
 * to swap out entire pages.
 *
 * If the table gets too big to fit in RAM:
 *    * Combine the Edens, Cloaks, and UnusedCargo dimensions.  Of the
 *      24 combinations, only 15 are legal: a 38% savings.
 *    * Reduce the size of the Fuel dimension to 3.  We only ever look
 *      backwards 2 units, so just rotate the logical values through
 *      the same 3 physical addresses.  This is good for an 82% savings.
 *    * Reduce the size of the Edens dimension from 3 to 2, for the
 *      same reasons as Fuel above.  33% savings.
 *    * Buy more ram.  (Just sayin'.  It's cheaper than you think.)
 *      
 */

// The official list of dimensions:
const (
	// Name                Num   Size  Description
	Edens       = iota //   1      3  # of Eden warp units (0 - 2 typically)
	Cloaks             //   2    1-2  # of Devices of Cloaking (0 or 1)
	UnusedCargo        //   3      4  # of unused cargo spaces (0 - 3 typically)
	Fuel               //   4     17  Hyper jump power left (0 - 16)
	Location           //   5     26  Location (which planet)
	Hold               //   6     15  Cargo bay contents (a *Commodity or nil)
	Traded             //   7      2  Traded yet?
	BuyFighters        //   8    1-2  Errand: Buy fighter drones
	BuyShields         //   9    1-2  Errand: Buy shield batteries
	Visit              //  10 1-2**N  Visit: Stop by these N planets in the route

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
	if *start_edens > eden_capacity {
		eden_capacity = *start_edens
	}
	cloak_capacity := bint(*cloak)
	dims := make([]int, NumDimensions)
	dims[Edens] = eden_capacity + 1
	dims[Cloaks] = cloak_capacity + 1
	dims[UnusedCargo] = eden_capacity + cloak_capacity + 1
	dims[Fuel] = *fuel + 1
	dims[Location] = len(data.Planets)
	dims[Hold] = len(data.Commodities) + 1
	dims[Traded] = 2
	dims[BuyFighters] = bint(*drones > 0) + 1
	dims[BuyShields] = bint(*batteries > 0) + 1
	dims[Visit] = 1 << uint(len(visit()))

	// Remind myself to add a line above when adding new dimensions
	for i, dim := range dims {
		if dim < 1 {
			panic(i)
		}
	}
	return dims
}

func StateTableSize(dims []int) int {
	product := 1
	for _, size := range dims {
		product *= size
	}
	return product
}

type State struct {
	value, from int32
}

const (
	FROM_ROOT = -2147483647 + iota
	FROM_UNINITIALIZED
	VALUE_UNINITIALIZED
	VALUE_BEING_EVALUATED
	VALUE_RUBISH
)

func EncodeIndex(dims, addr []int) int32 {
	index := addr[0]
	if addr[0] > dims[0] {
		panic(0)
	}
	for i := 1; i < NumDimensions; i++ {
		if addr[i] < 0 || addr[i] >= dims[i] {
			panic(i)
		}
		index = index*dims[i] + addr[i]
	}
	return int32(index)
}

func DecodeIndex(dims []int, index int32) []int {
	addr := make([]int, NumDimensions)
	for i := NumDimensions - 1; i > 0; i-- {
		addr[i] = int(index) % dims[i]
		index /= int32(dims[i])
	}
	addr[0] = int(index)
	return addr
}

func CreateStateTable(data planet_data, dims []int) []State {
	table := make([]State, StateTableSize(dims))
	for i := range table {
		table[i].value = VALUE_UNINITIALIZED
		table[i].from = FROM_UNINITIALIZED
	}

	addr := make([]int, NumDimensions)
	addr[Fuel] = *fuel
	addr[Edens] = *start_edens
	addr[Location] = data.p2i[*start]
	addr[Traded] = 1
	start_index := EncodeIndex(dims, addr)
	table[start_index].value = int32(*funds)
	table[start_index].from = FROM_ROOT

	return table
}

/* CellValue fills in the one cell at address addr by looking at all
 * the possible ways to reach this cell and selecting the best one. */

func Consider(data planet_data, dims []int, table []State, there []int, value_difference int, best_value *int32, best_source []int) {
	there_value := CellValue(data, dims, table, there)
	if value_difference < 0 && int32(-value_difference) > there_value {
		/* Can't afford this transition */
		return
	}
	possible_value := there_value + int32(value_difference)
	if possible_value > *best_value {
		*best_value = possible_value
		copy(best_source, there)
	}
}

var cell_filled_count int

func CellValue(data planet_data, dims []int, table []State, addr []int) int32 {
	my_index := EncodeIndex(dims, addr)
	if table[my_index].value == VALUE_BEING_EVALUATED {
		panic("Circular dependency")
	}
	if table[my_index].value != VALUE_UNINITIALIZED {
		return table[my_index].value
	}
	table[my_index].value = VALUE_BEING_EVALUATED

	best_value := int32(VALUE_RUBISH)
	best_source := make([]int, NumDimensions)
	other := make([]int, NumDimensions)
	copy(other, addr)
	planet := data.i2p[addr[Location]]

	/* Travel here */
	if addr[Traded] == 0 { /* Can't have traded immediately after traveling. */
		other[Traded] = 1 /* Travel from states that have done trading. */

		/* Travel here via a 2-fuel unit jump */
		if addr[Fuel]+2 < dims[Fuel] {
			other[Fuel] = addr[Fuel] + 2
			hole_index := (dims[Fuel] - 1) - (addr[Fuel] + 2)
			if hole_index >= len(flight_plan()) || addr[Location] != data.p2i[flight_plan()[hole_index]] {
				for other[Location] = 0; other[Location] < dims[Location]; other[Location]++ {
					if data.Planets[data.i2p[addr[Location]]].BeaconOn {
						Consider(data, dims, table, other, 0, &best_value, best_source)
					}
				}
			}
			other[Location] = addr[Location]
			other[Fuel] = addr[Fuel]
		}

		/* Travel here via a 1-fuel unit jump (a hyper hole) */
		if addr[Fuel]+1 < dims[Fuel] {
			hole_index := (dims[Fuel] - 1) - (addr[Fuel] + 1)
			if hole_index < len(flight_plan()) && addr[Location] == data.p2i[flight_plan()[hole_index]] {
				other[Fuel] = addr[Fuel] + 1
				for other[Location] = 0; other[Location] < dims[Location]; other[Location]++ {
					Consider(data, dims, table, other, 0, &best_value, best_source)
				}
				other[Location] = addr[Location]
				other[Fuel] = addr[Fuel]
			}
		}

		/* Travel here via Eden Warp Unit */
		if addr[Edens]+1 < dims[Edens] && addr[UnusedCargo] > 0 {
			_, available := data.Planets[data.i2p[addr[Location]]].RelativePrices["Eden Warp Units"]
			if !available {
				other[Edens] = addr[Edens] + 1
				if other[Hold] != 0 {
					other[UnusedCargo] = addr[UnusedCargo] - 1
				}
				for other[Location] = 0; other[Location] < dims[Location]; other[Location]++ {
					Consider(data, dims, table, other, 0, &best_value, best_source)
				}
				other[Location] = addr[Location]
				other[UnusedCargo] = addr[UnusedCargo]
				other[Edens] = addr[Edens]
			}
		}
		other[Traded] = addr[Traded]
	}

	/* Trade */
	if addr[Traded] == 1 {
		other[Traded] = 0

		/* Consider not trading */
		Consider(data, dims, table, other, 0, &best_value, best_source)

		if !data.Planets[data.i2p[addr[Location]]].Private {

			/* Sell */
			if addr[Hold] == 0 && addr[UnusedCargo] == 0 {
				for other[Hold] = 0; other[Hold] < dims[Hold]; other[Hold]++ {
					commodity := data.i2c[other[Hold]]
					if !data.Commodities[commodity].CanSell {
						continue
					}
					relative_price, available := data.Planets[planet].RelativePrices[commodity]
					if !available {
						// TODO: Dump cargo
						continue
					}
					base_price := data.Commodities[commodity].BasePrice
					absolute_price := float64(base_price) * float64(relative_price) / 100.0
					sell_price := int(absolute_price * 0.9)

					for other[UnusedCargo] = 0; other[UnusedCargo] < dims[UnusedCargo]; other[UnusedCargo]++ {
						quantity := *hold - (other[UnusedCargo] + other[Cloaks] + other[Edens])
						sale_value := quantity * sell_price
						Consider(data, dims, table, other, sale_value, &best_value, best_source)
					}
				}
				other[UnusedCargo] = addr[UnusedCargo]
				other[Hold] = addr[Hold]
			}

			/* Buy */
			other[Traded] = addr[Traded] /* Buy after selling */
			if addr[Hold] != 0 {
				commodity := data.i2c[addr[Hold]]
				if data.Commodities[commodity].CanSell {
					relative_price, available := data.Planets[planet].RelativePrices[commodity]
					if available {
						base_price := data.Commodities[commodity].BasePrice
						absolute_price := int(float64(base_price) * float64(relative_price) / 100.0)
						quantity := *hold - (addr[UnusedCargo] + addr[Cloaks] + addr[Edens])
						total_price := quantity * absolute_price
						other[Hold] = 0
						other[UnusedCargo] = 0
						Consider(data, dims, table, other, -total_price, &best_value, best_source)
						other[UnusedCargo] = addr[UnusedCargo]
						other[Hold] = addr[Hold]
					}
				}
			}
		}
		other[Traded] = addr[Traded]
	}

	/* Buy a Device of Cloaking */
	if addr[Cloaks] == 1 && addr[UnusedCargo] < dims[UnusedCargo]-1 {
		relative_price, available := data.Planets[data.i2p[addr[Location]]].RelativePrices["Device Of Cloakings"]
		if available {
			absolute_price := int(float64(data.Commodities["Device Of Cloakings"].BasePrice) * float64(relative_price) / 100.0)
			other[Cloaks] = 0
			if other[Hold] != 0 {
				other[UnusedCargo] = addr[UnusedCargo] + 1
			}
			Consider(data, dims, table, other, -absolute_price, &best_value, best_source)
			other[UnusedCargo] = addr[UnusedCargo]
			other[Cloaks] = addr[Cloaks]
		}
	}

	/* Buy Fighter Drones */
	if addr[BuyFighters] == 1 {
		relative_price, available := data.Planets[data.i2p[addr[Location]]].RelativePrices["Fighter Drones"]
		if available {
			absolute_price := int(float64(data.Commodities["Fighter Drones"].BasePrice) * float64(relative_price) / 100.0)
			other[BuyFighters] = 0
			Consider(data, dims, table, other, -absolute_price**drones, &best_value, best_source)
			other[BuyFighters] = addr[BuyFighters]
		}
	}

	/* Buy Shield Batteries */
	if addr[BuyShields] == 1 {
		relative_price, available := data.Planets[data.i2p[addr[Location]]].RelativePrices["Shield Batterys"]
		if available {
			absolute_price := int(float64(data.Commodities["Shield Batterys"].BasePrice) * float64(relative_price) / 100.0)
			other[BuyShields] = 0
			Consider(data, dims, table, other, -absolute_price**batteries, &best_value, best_source)
			other[BuyShields] = addr[BuyShields]
		}
	}

	/* Visit this planet */
	var i uint
	for i = 0; i < uint(len(visit())); i++ {
		if addr[Visit]&(1<<i) != 0 && visit()[i] == data.i2p[addr[Location]] {
			other[Visit] = addr[Visit] & ^(1 << i)
			Consider(data, dims, table, other, 0, &best_value, best_source)
		}
	}
	other[Visit] = addr[Visit]

	/* Buy Eden warp units */
	eden_limit := data.Commodities["Eden Warp Units"].Limit
	if addr[Edens] > 0 && addr[Edens] <= eden_limit {
		relative_price, available := data.Planets[data.i2p[addr[Location]]].RelativePrices["Eden Warp Units"]
		if available {
			absolute_price := int(float64(data.Commodities["Eden Warp Units"].BasePrice) * float64(relative_price) / 100.0)
			for quantity := addr[Edens]; quantity > 0; quantity-- {
				other[Edens] = addr[Edens] - quantity
				if addr[Hold] != 0 {
					other[UnusedCargo] = addr[UnusedCargo] + quantity
				}
				if other[UnusedCargo] < dims[UnusedCargo] {
					Consider(data, dims, table, other, -absolute_price*quantity, &best_value, best_source)
				}
			}
			other[Edens] = addr[Edens]
			other[UnusedCargo] = addr[UnusedCargo]
		}
	}

	// Check that we didn't lose track of any temporary modifications to other.
	for i := 0; i < NumDimensions; i++ {
		if addr[i] != other[i] {
			panic(i)
		}
	}

	// Sanity check: This cell was in state BEING_EVALUATED
	// the whole time that it was being evaluated.
	if table[my_index].value != VALUE_BEING_EVALUATED {
		panic(my_index)
	}

	// Record our findings
	table[my_index].value = best_value
	table[my_index].from = EncodeIndex(dims, best_source)

	// UI: Progress bar
	cell_filled_count++
	if cell_filled_count&0xff == 0 {
		print(fmt.Sprintf("\r%3.1f%%", 100*float64(cell_filled_count)/float64(StateTableSize(dims))))
	}

	return table[my_index].value
}

func FindBestState(data planet_data, dims []int, table []State) int32 {
	addr := make([]int, NumDimensions)
	addr[Edens] = *end_edens
	addr[Cloaks] = dims[Cloaks] - 1
	addr[BuyFighters] = dims[BuyFighters] - 1
	addr[BuyShields] = dims[BuyShields] - 1
	addr[Visit] = dims[Visit] - 1
	addr[Traded] = 1
	addr[Hold] = 0
	addr[UnusedCargo] = 0
	max_index := int32(-1)
	max_value := int32(0)
	max_fuel := 1
	if *fuel == 0 {
		max_fuel = 0
	}
	for addr[Fuel] = 0; addr[Fuel] <= max_fuel; addr[Fuel]++ {
		for addr[Location] = 0; addr[Location] < dims[Location]; addr[Location]++ {
			if len(end()) == 0 || end()[data.i2p[addr[Location]]] {
				index := EncodeIndex(dims, addr)
				value := CellValue(data, dims, table, addr)
				if value > max_value {
					max_value = value
					max_index = index
				}
			}
		}
	}
	return max_index
}

func Commas(n int32) (s string) {
	r := n % 1000
	n /= 1000
	for n > 0 {
		s = fmt.Sprintf(",%03d", r) + s
		r = n % 1000
		n /= 1000
	}
	s = fmt.Sprint(r) + s
	return
}

func DescribePath(data planet_data, dims []int, table []State, start int32) (description []string) {
	for index := start; table[index].from > FROM_ROOT; index = table[index].from {
		if table[index].from == FROM_UNINITIALIZED {
			panic(index)
		}
		var line string
		addr := DecodeIndex(dims, index)
		prev := DecodeIndex(dims, table[index].from)
		if addr[Fuel] != prev[Fuel] {
			from := data.i2p[prev[Location]]
			to := data.i2p[addr[Location]]
			line += fmt.Sprintf("Jump from %v to %v (%v hyper jump units)", from, to, prev[Fuel]-addr[Fuel])
		}
		if addr[Edens] == prev[Edens]-1 {
			from := data.i2p[prev[Location]]
			to := data.i2p[addr[Location]]
			line += fmt.Sprintf("Eden warp from %v to %v", from, to)
		}
		if addr[Hold] != prev[Hold] {
			if addr[Hold] == 0 {
				quantity := *hold - (prev[UnusedCargo] + prev[Edens] + prev[Cloaks])
				line += fmt.Sprintf("Sell %v %v", quantity, data.i2c[prev[Hold]])
			} else if prev[Hold] == 0 {
				quantity := *hold - (addr[UnusedCargo] + addr[Edens] + addr[Cloaks])
				line += fmt.Sprintf("Buy %v %v", quantity, data.i2c[addr[Hold]])
			} else {
				panic("Switched cargo?")
			}

		}
		if addr[Cloaks] == 1 && prev[Cloaks] == 0 {
			// TODO: Dump cloaks, convert from cargo?
			line += "Buy a Cloak"
		}
		if addr[Edens] > prev[Edens] {
			line += fmt.Sprint("Buy ", addr[Edens]-prev[Edens], " Eden Warp Units")
		}
		if addr[BuyShields] == 1 && prev[BuyShields] == 0 {
			line += fmt.Sprint("Buy ", *batteries, " Shield Batterys")
		}
		if addr[BuyFighters] == 1 && prev[BuyFighters] == 0 {
			line += fmt.Sprint("Buy ", *drones, " Fighter Drones")
		}
		if addr[Visit] != prev[Visit] {
			// TODO: verify that the bit chat changed is addr[Location]
			line += fmt.Sprint("Visit ", data.i2p[addr[Location]])
		}
		if line == "" && addr[Hold] == prev[Hold] && addr[Traded] != prev[Traded] {
			// The Traded dimension is for housekeeping.  It doesn't directly
			// correspond to in-game actions, so don't report transitions.
			continue
		}
		if line == "" {
			line = fmt.Sprint(prev, " -> ", addr)
		}
		description = append(description, fmt.Sprintf("%13v ", Commas(table[index].value))+line)
	}
	return
}

// (Example of a use case for generics in Go)
func IndexPlanets(m *map[string]Planet, start_at int) (map[string]int, []string) {
	e2i := make(map[string]int, len(*m)+start_at)
	i2e := make([]string, len(*m)+start_at)
	i := start_at
	for e := range *m {
		e2i[e] = i
		i2e[i] = e
		i++
	}
	return e2i, i2e
}
func IndexCommodities(m *map[string]Commodity, start_at int) (map[string]int, []string) {
	e2i := make(map[string]int, len(*m)+start_at)
	i2e := make([]string, len(*m)+start_at)
	i := start_at
	for e := range *m {
		e2i[e] = i
		i2e[i] = e
		i++
	}
	return e2i, i2e
}

func main() {
	flag.Parse()
	if *start == "" || *funds == 0 {
		print("--start and --funds are required.  --help for more\n")
		return
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	data := ReadData()
	if *drone_price > 0 {
		temp := data.Commodities["Fighter Drones"]
		temp.BasePrice = *drone_price
		data.Commodities["Fighter Drones"] = temp
	}
	if *battery_price > 0 {
		temp := data.Commodities["Shield Batterys"]
		temp.BasePrice = *battery_price
		data.Commodities["Shield Batterys"] = temp
	}
	data.p2i, data.i2p = IndexPlanets(&data.Planets, 0)
	data.c2i, data.i2c = IndexCommodities(&data.Commodities, 1)
	dims := DimensionSizes(data)
	table := CreateStateTable(data, dims)
	best := FindBestState(data, dims, table)
	print("\n")
	if best == -1 {
		print("Cannot acheive success criteria\n")
	} else {
		description := DescribePath(data, dims, table, best)
		for i := len(description) - 1; i >= 0; i-- {
			fmt.Println(description[i])
		}
	}
}
