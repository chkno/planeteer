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
	Name      string
	BasePrice int
	CanSell   bool
	Limit     int
}

type planet_data struct {
	Commodities []Commodity
	Planets []struct {
		Name     string
		BeaconOn bool
		/* Use relative prices rather than absolute prices because you
		   can get relative prices without traveling to each planet. */
		RelativePrices map [string] int
	}
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

func TradeValue(data planet_data,
                from_index, to_index, commodity_index, quantity int) int {

	commodity := &data.Commodities[commodity_index]
	if !commodity.CanSell {
		return 0
	}

	from_planet := &data.Planets[from_index]
	from_relative_price, from_available := from_planet.RelativePrices[commodity.Name]
	if !from_available {
		return 0
	}

	to_planet := &data.Planets[to_index]
	to_relative_price, to_available := to_planet.RelativePrices[commodity.Name]
	if !to_available {
		return 0
	}

	from_absolute_price := from_relative_price * commodity.BasePrice
	to_absolute_price := to_relative_price * commodity.BasePrice
	buy_price := from_absolute_price
	sell_price := int(float64(to_absolute_price) * 0.9)
	return (sell_price - buy_price) * quantity

}

func FindBestTrades(data planet_data) [][]*Commodity {
	best := make([][]*Commodity, len(data.Planets))
	for from_index := range data.Planets {
		best[from_index] = make([]*Commodity, len(data.Planets))
		for to_index := range data.Planets {
			best_gain := 0
			for commodity_index := range data.Commodities {
				gain := TradeValue(data, from_index, to_index, commodity_index, 1)
				if gain > best_gain {
					best[from_index][to_index] = &data.Commodities[commodity_index]
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
			if best_trades[from_index][to_index] != nil {
				best_trade = best_trades[from_index][to_index].Name
			}
			fmt.Printf("%s to %s: %s\n", from_planet.Name, to_planet.Name, best_trade)
		}
	}
}
