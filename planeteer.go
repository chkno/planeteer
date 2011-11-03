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

var funds = flag.Int("funds", 0,
	"Starting funds")

var start = flag.String("start", "",
	"The planet to start at")

var flight_plan_string = flag.String("flight_plan", "",
	"Your hidey-holes for the day, comma-separated.")

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
	if *visit_string == "" {
		return []string{}
	}
	return strings.Split(*visit_string, ",")
}

func flight_plan() []string {
	if *flight_plan_string == "" {
		return []string{}
	}
	return strings.Split(*flight_plan_string, ",")
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
 * cache lines, and if they are large enough, prevent the memory manager
 * from allocating pages for these areas at all.
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
	dims[NeedFighters] = bint(*drones > 0) + 1
	dims[NeedShields] = bint(*batteries > 0) + 1
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
	value, from int
}

func EncodeIndex(dims, addr []int) int {
	index := addr[0]
	if addr[0] > dims[0] {
		panic(0)
	}
	for i := 1; i < NumDimensions; i++ {
		if addr[i] > dims[i] {
			panic(i)
		}
		index = index*dims[i] + addr[i]
	}
	return index
}

func DecodeIndex(dims []int, index int) []int {
	addr := make([]int, NumDimensions)
	for i := NumDimensions - 1; i > 0; i-- {
		addr[i] = index % dims[i]
		index /= dims[i]
	}
	addr[0] = index
	return addr
}

func InitializeStateTable(data planet_data, dims []int) []State {
	table := make([]State, StateTableSize(dims))

	addr := make([]int, NumDimensions)
	addr[Fuel] = *fuel
	addr[Edens] = *start_edens
	addr[Location] = data.p2i[*start]
	table[EncodeIndex(dims, addr)].value = *funds

	return table
}

/* These four fill procedures fill in the cell at address addr by
 * looking at all the possible ways to reach this cell and selecting
 * the best one.
 *
 * The other obvious implementation choice is to do this the other way
 * around -- for each cell, conditionally overwrite all the other cells
 * that are reachable *from* the considered cell.  We choose gathering
 * reads over scattering writes to avoid having to take a bunch of locks.
 */

func UpdateCell(table []State, here, there, value_difference int) {
	possible_value := table[there].value + value_difference
	if table[there].value > 0 && possible_value > table[here].value {
		table[here].value = possible_value
		table[here].from = there
	}
}

func FillCellByArriving(data planet_data, dims []int, table []State, addr []int) {
	my_index := EncodeIndex(dims, addr)
	other := make([]int, NumDimensions)
	copy(other, addr)

	/* Travel here via a 2-fuel unit jump */
	if addr[Fuel]+2 < dims[Fuel] {
		other[Fuel] = addr[Fuel] + 2
		for other[Location] = 0; other[Location] < dims[Location]; other[Location]++ {
			UpdateCell(table, my_index, EncodeIndex(dims, other), 0)
		}
		other[Location] = addr[Location]
		other[Fuel] = addr[Fuel]
	}

	/* Travel here via a hidey hole */
	if addr[Fuel]+1 < dims[Fuel] {
		hole_index := (dims[Fuel] - 1) - (addr[Fuel] + 1)
		if hole_index < len(flight_plan()) && addr[Location] == data.p2i[flight_plan()[hole_index]] {
			other[Fuel] = addr[Fuel] + 1
			for other[Location] = 0; other[Location] < dims[Location]; other[Location]++ {
				UpdateCell(table, my_index, EncodeIndex(dims, other), 0)
			}
			other[Location] = addr[Location]
			other[Fuel] = addr[Fuel]
		}
	}

	/* Travel here via Eden Warp Unit */
	for other[Edens] = addr[Edens] + 1; other[Edens] < dims[Edens]; other[Edens]++ {
		for other[Location] = 0; other[Location] < dims[Location]; other[Location]++ {
			UpdateCell(table, my_index, EncodeIndex(dims, other), 0)
		}
	}
	other[Location] = addr[Location]
	other[Edens] = addr[Edens]
}

func FillCellBySelling(data planet_data, dims []int, table []State, addr []int) {
	if addr[Hold] > 0 {
		// Can't sell and still have cargo
		return
	}
	if addr[UnusedCargo] > 0 {
		// Can't sell everything and still have 'unused' holds
		return
	}
	my_index := EncodeIndex(dims, addr)
	other := make([]int, NumDimensions)
	copy(other, addr)
	planet := data.i2p[addr[Location]]
	for other[Hold] = 0; other[Hold] < dims[Hold]; other[Hold]++ {
		commodity := data.i2c[other[Hold]]
		if !data.Commodities[commodity].CanSell {
			// TODO: Dump cargo
			continue
		}
		relative_price, available := data.Planets[planet].RelativePrices[commodity]
		if !available {
			continue
		}
		base_price := data.Commodities[commodity].BasePrice
		absolute_price := float64(base_price) * float64(relative_price) / 100.0
		sell_price := int(absolute_price * 0.9)

		for other[UnusedCargo] = 0; other[UnusedCargo] < dims[UnusedCargo]; other[UnusedCargo]++ {

			quantity := *hold - (other[UnusedCargo] + other[Cloaks] + other[Edens])
			sale_value := quantity * sell_price
			UpdateCell(table, my_index, EncodeIndex(dims, other), sale_value)
		}
	}
	other[UnusedCargo] = addr[UnusedCargo]
}

func FillCellByBuying(data planet_data, dims []int, table []State, addr []int) {
	if addr[Hold] == 0 {
		// Can't buy and then have nothing
		return
	}
	my_index := EncodeIndex(dims, addr)
	other := make([]int, NumDimensions)
	copy(other, addr)
	planet := data.i2p[addr[Location]]
	commodity := data.i2c[addr[Hold]]
	if !data.Commodities[commodity].CanSell {
		return
	}
	relative_price, available := data.Planets[planet].RelativePrices[commodity]
	if !available {
		return
	}
	base_price := data.Commodities[commodity].BasePrice
	absolute_price := int(float64(base_price) * float64(relative_price) / 100.0)
	quantity := *hold - (addr[UnusedCargo] + addr[Cloaks] + addr[Edens])
	total_price := quantity * absolute_price
	other[Hold] = 0
	other[UnusedCargo] = 0
	UpdateCell(table, my_index, EncodeIndex(dims, other), -total_price)
	other[UnusedCargo] = addr[UnusedCargo]
	other[Hold] = addr[Hold]
}

func FillCellByMisc(data planet_data, dims []int, table []State, addr []int) {
	my_index := EncodeIndex(dims, addr)
	other := make([]int, NumDimensions)
	copy(other, addr)
	/* Buy Eden warp units */
	/* Buy a Device of Cloaking */
	if addr[Cloaks] == 1 && addr[UnusedCargo] < dims[UnusedCargo]-1 {
		relative_price, available := data.Planets[data.i2p[addr[Location]]].RelativePrices["Device Of Cloakings"]
		if available {
			absolute_price := int(float64(data.Commodities["Device Of Cloakings"].BasePrice) * float64(relative_price) / 100.0)
			other[Cloaks] = 0
			if other[Hold] != 0 {
				other[UnusedCargo] = addr[UnusedCargo] + 1
			}
			UpdateCell(table, my_index, EncodeIndex(dims, other), -absolute_price)
			other[UnusedCargo] = addr[UnusedCargo]
			other[Cloaks] = addr[Cloaks]
		}
	}
	/* Silly: Dump a Device of Cloaking */
	/* Buy Fighter Drones */
	/* Buy Shield Batteries */
	/* Visit this planet */
}

func FillStateTable2Iteration(data planet_data, dims []int, table []State,
addr []int, f func(planet_data, []int, []State, []int)) {
	/* TODO: Justify the safety of the combination of this dimension
	 * iteration and the various phases f.  */
	for addr[Hold] = 0; addr[Hold] < dims[Hold]; addr[Hold]++ {
		for addr[Cloaks] = 0; addr[Cloaks] < dims[Cloaks]; addr[Cloaks]++ {
			for addr[UnusedCargo] = 0; addr[UnusedCargo] < dims[UnusedCargo]; addr[UnusedCargo]++ {
				for addr[NeedFighters] = 0; addr[NeedFighters] < dims[NeedFighters]; addr[NeedFighters]++ {
					for addr[NeedShields] = 0; addr[NeedShields] < dims[NeedShields]; addr[NeedShields]++ {
						for addr[Visit] = 0; addr[Visit] < dims[Visit]; addr[Visit]++ {
							f(data, dims, table, addr)
						}
					}
				}
			}
		}
	}
}

func FillStateTable2(data planet_data, dims []int, table []State,
fuel_remaining, edens_remaining int, planet string, barrier chan<- bool) {
	addr := make([]int, len(dims))
	addr[Edens] = edens_remaining
	addr[Fuel] = fuel_remaining
	addr[Location] = data.p2i[planet]
	FillStateTable2Iteration(data, dims, table, addr, FillCellByArriving)
	FillStateTable2Iteration(data, dims, table, addr, FillCellBySelling)
	FillStateTable2Iteration(data, dims, table, addr, FillCellByBuying)
	FillStateTable2Iteration(data, dims, table, addr, FillCellByMisc)
	barrier <- true
}

/* Filling the state table is a set of nested for loops NumDimensions deep.
 * We split this into two procedures: 1 and 2.  #1 is the outer, slowest-
 * changing indexes.  #1 fires off many calls to #2 that run in parallel.
 * The order of the nesting of the dimensions, the order of iteration within
 * each dimension, and where the 1 / 2 split is placed are carefully chosen 
 * to make this arrangement safe.
 *
 * Outermost two layers: Go from high-energy states (lots of fuel, edens) to
 * low-energy state.  These must be processed sequentially and in this order
 * because you travel through high-energy states to get to the low-energy
 * states.
 *
 * Third layer: Planet.  This is a good layer to parallelize on.  There's
 * high enough cardinality that we don't have to mess with parallelizing
 * multiple layers for good utilization (on 2011 machines).  Each thread
 * works on one planet's states and need not synchronize with peer threads.
 */
func FillStateTable1(data planet_data, dims []int, table []State) {
	barrier := make(chan bool, len(data.Planets))
	eden_capacity := data.Commodities["Eden Warp Units"].Limit
	work_units := (float64(*fuel) + 1) * (float64(eden_capacity) + 1)
	work_done := 0.0
	for fuel_remaining := *fuel; fuel_remaining >= 0; fuel_remaining-- {
		for edens_remaining := eden_capacity; edens_remaining >= 0; edens_remaining-- {
			for planet := range data.Planets {
				go FillStateTable2(data, dims, table, fuel_remaining,
					edens_remaining, planet, barrier)
			}
			for _ = range data.Planets {
				<-barrier
			}
			work_done++
			print(fmt.Sprintf("\r%3.0f%%", 100*work_done/work_units))
		}
	}
	print("\n")
}

func FindBestState(data planet_data, dims []int, table []State) int {
	addr := make([]int, NumDimensions)
	addr[Edens] = *end_edens
	addr[Cloaks] = dims[Cloaks] - 1
	addr[NeedFighters] = dims[NeedFighters] - 1
	addr[NeedShields] = dims[NeedShields] - 1
	addr[Visit] = dims[Visit] - 1
	// Fuel, Hold, UnusedCargo left at 0
	max_index := -1
	max_value := 0
	for addr[Location] = 0; addr[Location] < dims[Location]; addr[Location]++ {
		index := EncodeIndex(dims, addr)
		if table[index].value > max_value {
			max_value = table[index].value
			max_index = index
		}
	}
	return max_index
}

func Commas(n int) (s string) {
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

func DescribePath(data planet_data, dims []int, table []State, start int) (description []string) {
	for index := start; index > 0 && table[index].from > 0; index = table[index].from {
		line := fmt.Sprintf("%13v", Commas(table[index].value))
		addr := DecodeIndex(dims, index)
		prev := DecodeIndex(dims, table[index].from)
		if addr[Location] != prev[Location] {
			from := data.i2p[prev[Location]]
			to := data.i2p[addr[Location]]
			if addr[Fuel] != prev[Fuel] {
				line += fmt.Sprintf(" Jump from %v to %v (%v reactor units)", from, to, prev[Fuel]-addr[Fuel])
			} else if addr[Edens] != prev[Edens] {
				line += fmt.Sprintf(" Eden warp from %v to %v", from, to)
			} else {
				panic("Traveling without fuel?")
			}
		}
		if addr[Hold] != prev[Hold] {
			if addr[Hold] == 0 {
				quantity := *hold - (prev[UnusedCargo] + prev[Edens] + prev[Cloaks])
				line += fmt.Sprintf(" Sell %v %v", quantity, data.i2c[prev[Hold]])
			} else if prev[Hold] == 0 {
				quantity := *hold - (addr[UnusedCargo] + addr[Edens] + addr[Cloaks])
				line += fmt.Sprintf(" Buy %v %v", quantity, data.i2c[addr[Hold]])
			} else {
				panic("Switched cargo?")
			}

		}
		if addr[Cloaks] == 1 && prev[Cloaks] == 0 {
			// TODO: Dump cloaks, convert from cargo?
			line += " Buy a Cloak"
		}
		description = append(description, line)
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
	data := ReadData()
	data.p2i, data.i2p = IndexPlanets(&data.Planets, 0)
	data.c2i, data.i2c = IndexCommodities(&data.Commodities, 1)
	dims := DimensionSizes(data)
	table := InitializeStateTable(data, dims)
	FillStateTable1(data, dims, table)
	best := FindBestState(data, dims, table)
	if best == -1 {
		print("Cannot acheive success criteria\n")
	} else {
		description := DescribePath(data, dims, table, best)
		for i := len(description) - 1; i >= 0; i-- {
			fmt.Println(description[i])
		}
	}
}
