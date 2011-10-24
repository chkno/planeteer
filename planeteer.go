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
import "json"
import "os"
import "fmt"

var datafile = flag.String("planet_data_file", "planet-data",
	"The file to read planet data from")

type Commodity struct {
	BasePrice int
	CanSell   bool
	Limit     int
}
type Planet struct {
	Name     string
	BeaconOn bool
	/* Use relative prices rather than absolute prices because you
	   can get relative prices without traveling to each planet. */
	RelativePrices map [string] int
}
type planet_data struct {
	Commodities map [string] Commodity
	Planets []Planet
}

func ReadData() (data planet_data) {
	f, err := os.Open(*datafile)
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

/* What is the value of hauling 'commodity' from 'from' to 'to'?
 * Take into account the available funds and the available cargo space. */
func TradeValue(data planet_data,
                from, to *Planet,
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
	best := make([][]string, len(data.Planets))
	for from_index, from_planet := range data.Planets {
		best[from_index] = make([]string, len(data.Planets))
		for to_index, to_planet := range data.Planets {
			best_gain := 0
			price_list := from_planet.RelativePrices
			if len(to_planet.RelativePrices) < len(from_planet.RelativePrices) {
				price_list = to_planet.RelativePrices
			}
			for commodity := range price_list {
				gain := TradeValue(data,
				                   &from_planet,
				                   &to_planet,
				                   commodity,
				                   10000000,
				                   1)
				if gain > best_gain {
					best[from_index][to_index] = commodity
					gain = best_gain
				}
			}
		}
	}
	return best
}

func main() {
	flag.Parse()
	data := ReadData()
	best_trades := FindBestTrades(data)
	for from_index, from_planet := range data.Planets {
		for to_index, to_planet := range data.Planets {
			best_trade := "(nothing)"
			if best_trades[from_index][to_index] != "" {
				best_trade = best_trades[from_index][to_index]
			}
			fmt.Printf("%s to %s: %s\n", from_planet.Name, to_planet.Name, best_trade)
		}
	}
}
