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

type planet_data struct {
	Commodities []struct {
		Name      string
		BasePrice int
		CanSell   bool
		Limit     int
	}
	Planets []struct {
		Name     string
		BeaconOn bool
		/* Use relative prices rather than absolute prices because you
		   can get relative prices without traveling to each planet. */
		RelativePrices []struct {
			Name  string
			Value int
		}
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

func main() {
	flag.Parse()
	data := ReadData()
	fmt.Printf("%v", data)
}
