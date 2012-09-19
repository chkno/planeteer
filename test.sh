#!/bin/bash

exit_status=0

function t() {
  expected=$1
  shift
  diff -u --label Expected <( echo -n "$expected" ) \
          --label Actual <( ./planeteer "$@" ) || exit_status=1
}

t '      850,000 Buy 300 Medical Units
      850,000 Jump from Earth to Sickonia (2 hyper jump units)
    3,220,600 Sell 300 Medical Units
    2,770,600 Buy 300 Heating Units
    2,770,600 Jump from Sickonia to Hothor (2 hyper jump units)
    5,200,600 Sell 300 Heating Units
    4,999,600 Buy 300 Medical Units
    4,999,600 Jump from Hothor to Sickonia (2 hyper jump units)
    7,370,200 Sell 300 Medical Units
    6,920,200 Buy 300 Heating Units
    6,920,200 Jump from Sickonia to Hothor (2 hyper jump units)
    9,350,200 Sell 300 Heating Units
    9,149,200 Buy 300 Medical Units
    9,149,200 Jump from Hothor to Sickonia (2 hyper jump units)
   11,519,800 Sell 300 Medical Units
   11,069,800 Buy 300 Heating Units
   11,069,800 Jump from Sickonia to Hothor (2 hyper jump units)
   13,499,800 Sell 300 Heating Units
   13,298,800 Buy 300 Medical Units
   13,298,800 Jump from Hothor to Sickonia (2 hyper jump units)
   15,669,400 Sell 300 Medical Units
   15,219,400 Buy 300 Heating Units
   15,219,400 Jump from Sickonia to Hothor (2 hyper jump units)
   17,649,400 Sell 300 Heating Units
' --funds 1000000 --start Earth

t '' --funds 1000000 --start Earth --fuel 0

t '      850,000 Buy 300 Medical Units
      850,000 Jump from Earth to HomeWorld (2 hyper jump units)
    3,256,900 Sell 300 Medical Units
' --funds 1000000 --start Earth --fuel 2

t '      850,000 Buy 300 Medical Units
      850,000 Jump from Earth to Sickonia (2 hyper jump units)
    3,220,600 Sell 300 Medical Units
    2,770,600 Buy 300 Heating Units
    2,770,600 Jump from Sickonia to Hothor (2 hyper jump units)
    5,200,600 Sell 300 Heating Units
' --funds 1000000 --start Earth --fuel 4

t '      850,000 Buy 300 Medical Units
      850,000 Jump from Earth to Sickonia (2 hyper jump units)
    3,220,600 Sell 300 Medical Units
    2,770,600 Buy 300 Heating Units
    2,770,600 Jump from Sickonia to Hothor (2 hyper jump units)
    5,200,600 Sell 300 Heating Units
    4,999,600 Buy 300 Medical Units
    4,999,600 Jump from Hothor to HomeWorld (2 hyper jump units)
    7,406,500 Sell 300 Medical Units
' --funds 1000000 --start Earth --fuel 6

t '    3,900,000 Buy 300 AntiCloak Scanners
    3,900,000 Jump from Metallica to Loony (2 hyper jump units)
   32,790,000 Sell 300 AntiCloak Scanners
' --funds 30000000 --start Metallica --fuel 2

t '   19,934,000 Buy 300 Ground Weapons
   19,934,000 Jump from Metallica to Tribonia (2 hyper jump units)
   21,779,300 Sell 300 Ground Weapons
' --funds 20000000 --start Metallica --fuel 2

t '      850,000 Buy 300 Medical Units
      850,000 Jump from Earth to Sickonia (2 hyper jump units)
    3,220,600 Sell 300 Medical Units
    2,770,600 Buy 300 Heating Units
    2,770,600 Jump from Sickonia to Hothor (2 hyper jump units)
    5,200,600 Sell 300 Heating Units
    4,999,600 Buy 300 Medical Units
    4,999,600 Jump from Hothor to HomeWorld (2 hyper jump units)
    7,406,500 Sell 300 Medical Units
    7,303,000 Buy 300 Novelty Packs
    7,303,000 Jump from HomeWorld to Gojuon (2 hyper jump units)
    8,486,500 Sell 300 Novelty Packs

'$'\r''    863,700 Cost of --end Gojuon
' --funds 1000000 --start Earth --end Gojuon --fuel 8

t '      996,500 Buy 7 Medical Units
      996,500 Jump from Earth to Sickonia (2 hyper jump units)
    1,051,814 Sell 7 Medical Units
    1,041,314 Buy 7 Heating Units
    1,041,314 Jump from Sickonia to Hothor (2 hyper jump units)
    1,098,014 Sell 7 Heating Units
' --funds 1000000 --start Earth --fuel 4 --hold 7

t '    1,121,500 Sell 300 Clothes Bundles
      971,500 Buy 300 Medical Units
      971,500 Jump from Earth to HomeWorld (2 hyper jump units)
    3,378,400 Sell 300 Medical Units
' --funds 1000000 --start Earth --fuel 2 --start_hold 'Clothes Bundles'

t '    1,000,000 Jump from Earth to Eden (1 hyper jump units)
      950,000 Buy 2 Eden Warp Units
      923,776 Buy 298 Tree Growth Kits
      923,776 Eden warp from Eden to Zoolie
    3,026,464 Sell 298 Tree Growth Kits
    2,909,854 Buy 299 Medical Units
    2,909,854 Eden warp from Zoolie to HomeWorld
    5,308,731 Sell 299 Medical Units
    5,308,731 Jump from HomeWorld to Eden (1 hyper jump units)
    5,258,731 Buy 2 Eden Warp Units
    5,232,507 Buy 298 Tree Growth Kits
    5,232,507 Eden warp from Eden to Zoolie
    7,335,195 Sell 298 Tree Growth Kits
    7,218,585 Buy 299 Medical Units
    7,218,585 Eden warp from Zoolie to HomeWorld
    9,617,462 Sell 299 Medical Units
    9,617,462 Jump from HomeWorld to Eden (1 hyper jump units)
    9,567,462 Buy 2 Eden Warp Units
    9,541,238 Buy 298 Tree Growth Kits
    9,541,238 Eden warp from Eden to Zoolie
   11,643,926 Sell 298 Tree Growth Kits
   11,527,316 Buy 299 Medical Units
   11,527,316 Eden warp from Zoolie to HomeWorld
   13,926,193 Sell 299 Medical Units
' --funds 1000000 --start Earth --fuel 3 --flight_plan Eden,Eden,Eden

t '      965,500 Buy 300 Medical Units
      965,500 Jump from Medoca to Eden (1 hyper jump units)
    1,100,500 Sell 300 Medical Units
    1,050,500 Buy 2 Eden Warp Units

'$'\r''Use 1 extra edens, make an extra 2,256,084 ( 2,256,084 per eden)
'$'\r''Use 2 extra edens, make an extra 4,358,731 ( 2,179,365 per eden)
' --funds 1000000 --start Medoca --fuel 1 --flight_plan Eden --end_edens 2

t '      965,500 Buy 300 Medical Units
      965,500 Jump from Medoca to Eden (1 hyper jump units)
    1,100,500 Sell 300 Medical Units
    1,050,500 Buy 2 Eden Warp Units
' --funds 1000000 --start Medoca --fuel 1 --flight_plan Eden --end_edens 2 --extra_stats=false

t '      907,322 Buy 298 Ground Weapons
      907,322 Eden warp from Desha Rockna to Tribonia
    2,740,320 Sell 298 Ground Weapons
    2,558,528 Buy 299 Tree Growth Kits
    2,558,528 Jump from Tribonia to Zoolie (1 hyper jump units)
    4,668,272 Sell 299 Tree Growth Kits
    4,551,662 Buy 299 Medical Units
    4,551,662 Eden warp from Zoolie to Sickonia
    6,914,360 Sell 299 Medical Units
    6,464,360 Buy 300 Heating Units
    6,464,360 Jump from Sickonia to Hothor (1 hyper jump units)
    8,894,360 Sell 300 Heating Units
' --funds 1000000 --start 'Desha Rockna' --start_edens 2 --fuel 2 --flight_plan Zoolie,Hothor

t '      550,000 Buy 300 Heating Units
      550,000 Jump from Dune to Hothor (2 hyper jump units)
    2,980,000 Sell 300 Heating Units
    2,974,400 Buy a Cloak
    2,774,070 Buy 299 Medical Units
    2,774,070 Jump from Hothor to HomeWorld (2 hyper jump units)
    5,172,947 Sell 299 Medical Units
' --funds 1000000 --start Dune --fuel 4 --cloak

t '  173,900,000 Buy 300 AntiCloak Scanners
  173,900,000 Jump from Metallica to Loony (2 hyper jump units)
  202,790,000 Sell 300 AntiCloak Scanners
  201,965,000 Buy 300 Device Of Cloakings
  201,965,000 Jump from Loony to WeaponWorld (2 hyper jump units)
   76,335,368 Buy 654321 Fighter Drones
   77,658,368 Sell 300 Device Of Cloakings

'$'\r''Drones were 194.26 each
' --funds 200000000 --start Metallica --fuel 4 --drones 654321

t '   29,786,400 Buy 300 Tree Growth Kits
   29,786,400 Jump from Gojuon to Metallica (2 hyper jump units)
   23,913,582 Buy 87654 Shield Batterys
   25,157,682 Sell 300 Tree Growth Kits

'$'\r''Batteries were 77.88 each
' --funds 30000000 --start Gojuon --fuel 3 --batteries 87654

t '   29,908,800 Buy 300 Ground Weapons
   29,908,800 Jump from Gojuon to Tribonia (2 hyper jump units)
   31,754,100 Sell 300 Ground Weapons
   31,571,700 Buy 300 Tree Growth Kits
   31,571,700 Jump from Tribonia to Metallica (2 hyper jump units)
   25,698,882 Buy 87654 Shield Batterys
   26,942,982 Sell 300 Tree Growth Kits

'$'\r''Batteries were 82.71 each
' --funds 30000000 --start Gojuon --fuel 4 --batteries 87654

t '   27,808,650 Buy 87654 Shield Batterys
   27,358,650 Buy 300 Heating Units
   27,358,650 Jump from Dune to Hothor (2 hyper jump units)
   29,788,650 Sell 300 Heating Units
   29,587,650 Buy 300 Medical Units
   29,587,650 Jump from Hothor to HomeWorld (2 hyper jump units)
   31,994,550 Sell 300 Medical Units

'$'\r''Batteries were 25.00 each
' --funds 30000000 --start Dune --fuel 4 --batteries 87654 --battery_price 18

t '   29,472,000 Buy 300 Jewels
   29,472,000 Jump from Dune to WeaponWorld (2 hyper jump units)
   27,806,574 Buy 87654 Shield Batterys
   29,356,374 Sell 300 Jewels
    2,956,374 Buy 300 AntiCloak Scanners
    2,956,374 Jump from WeaponWorld to Loony (2 hyper jump units)
   31,846,374 Sell 300 AntiCloak Scanners

'$'\r''Batteries were 26.69 each
' --funds 30000000 --start Dune --fuel 4 --batteries 87654 --battery_price 19

t '   29,472,000 Buy 300 Jewels
   29,472,000 Jump from Dune to WeaponWorld (2 hyper jump units)
   27,368,304 Buy 87654 Shield Batterys
   28,918,104 Sell 300 Jewels
    2,518,104 Buy 300 AntiCloak Scanners
    2,518,104 Jump from WeaponWorld to Loony (2 hyper jump units)
   31,408,104 Sell 300 AntiCloak Scanners

'$'\r''Batteries were 31.69 each
' --funds 30000000 --start Dune --fuel 4 --batteries 87654 --battery_price 24

t '   30,000,000 Jump from Dune to Metallica (2 hyper jump units)
   28,597,536 Buy 87654 Shield Batterys
    2,497,536 Buy 300 AntiCloak Scanners
    2,497,536 Jump from Metallica to Loony (2 hyper jump units)
   31,387,536 Sell 300 AntiCloak Scanners

'$'\r''Batteries were 31.93 each
' --funds 30000000 --start Dune --fuel 4 --batteries 87654 --battery_price 25

t '   30,000,000 Jump from Dune to Metallica (2 hyper jump units)
   26,055,570 Buy 87654 Shield Batterys
       42,570 Buy 299 AntiCloak Scanners
       42,570 Jump from Metallica to Loony (2 hyper jump units)
   28,836,270 Sell 299 AntiCloak Scanners

'$'\r''Batteries were 61.03 each
' --funds 30000000 --start Dune --fuel 4 --batteries 87654 --battery_price 68

t '   29,871,900 Buy 300 Ground Weapons
   29,871,900 Jump from Dune to Tribonia (2 hyper jump units)
   31,717,200 Sell 300 Ground Weapons
   31,534,800 Buy 300 Tree Growth Kits
   31,534,800 Jump from Tribonia to Metallica (2 hyper jump units)
   27,502,716 Buy 87654 Shield Batterys
   28,746,816 Sell 300 Tree Growth Kits

'$'\r''Batteries were 62.05 each
' --funds 30000000 --start Dune --fuel 4 --batteries 87654 --battery_price 69

t '  137,826,324 Buy 57423 Fighter Drones
  111,126,324 Buy 300 AntiCloak Scanners
  111,126,324 Jump from Norhaven to Loony (2 hyper jump units)
  140,016,324 Sell 300 AntiCloak Scanners
  139,566,324 Buy 300 Heating Units
  139,566,324 Jump from Loony to Hothor (2 hyper jump units)
  141,996,324 Sell 300 Heating Units

'$'\r''Drones were 212.00 each
' --funds 150000000 --start Norhaven --fuel 4 --drones 57423
t '  149,802,000 Buy 300 Medical Units
  149,802,000 Jump from Norhaven to HomeWorld (2 hyper jump units)
  152,208,900 Sell 300 Medical Units
  151,470,900 Buy 300 Jewels
  151,470,900 Jump from HomeWorld to WeaponWorld (2 hyper jump units)
  140,445,492 Buy 57424 Fighter Drones
  141,995,292 Sell 300 Jewels

'$'\r''Drones were 212.01 each
' --funds 150000000 --start Norhaven --fuel 4 --drones 57424

t '  137,306,550 Buy 60445 Fighter Drones
  110,606,550 Buy 300 AntiCloak Scanners
  110,606,550 Jump from Norhaven to Loony (2 hyper jump units)
  139,496,550 Sell 300 AntiCloak Scanners
  139,046,550 Buy 300 Heating Units
  139,046,550 Jump from Loony to Hothor (2 hyper jump units)
  141,476,550 Sell 300 Heating Units

'$'\r''Drones were 210.00 each
' --funds 150000000 --start Norhaven --fuel 4 --drones 60445 --drone_price 199

t '  149,802,000 Buy 300 Medical Units
  149,802,000 Jump from Norhaven to HomeWorld (2 hyper jump units)
  152,208,900 Sell 300 Medical Units
  151,470,900 Buy 300 Jewels
  151,470,900 Jump from HomeWorld to WeaponWorld (2 hyper jump units)
  139,925,714 Buy 60446 Fighter Drones
  141,475,514 Sell 300 Jewels

'$'\r''Drones were 210.01 each
' --funds 150000000 --start Norhaven --fuel 4 --drones 60446 --drone_price 199

t '   98,500,000 Buy 300 Device Of Cloakings
   98,500,000 Jump from Earth to Volcana (1 hyper jump units)
  101,024,500 Sell 300 Device Of Cloakings
  100,943,500 Buy 300 Heating Units
  100,943,500 Jump from Volcana to Richiana (1 hyper jump units)
  103,227,700 Sell 300 Heating Units
  103,227,700 Jump from Richiana to Eden (1 hyper jump units)
  103,177,700 Buy 2 Eden Warp Units
  103,028,700 Buy 298 Medical Units
  103,028,700 Eden warp from Eden to HomeWorld
  105,419,554 Sell 298 Medical Units
  105,383,973 Buy 299 Ground Weapons
  105,383,973 Jump from HomeWorld to Tribonia (1 hyper jump units)
  107,223,122 Sell 299 Ground Weapons
  107,041,330 Buy 299 Tree Growth Kits
  107,041,330 Jump from Tribonia to Medoca (1 hyper jump units)
  108,046,568 Sell 299 Tree Growth Kits
  108,012,183 Buy 299 Medical Units
  108,012,183 Jump from Medoca to Sickonia (1 hyper jump units)
  110,374,881 Sell 299 Medical Units
  109,926,381 Buy 299 Heating Units
  109,926,381 Eden warp from Sickonia to Hothor
  112,348,281 Sell 299 Heating Units
  112,147,281 Buy 300 Medical Units
  112,147,281 Jump from Hothor to Sickonia (2 hyper jump units)
  114,517,881 Sell 300 Medical Units
  114,067,881 Buy 300 Heating Units
  114,067,881 Jump from Sickonia to Hothor (2 hyper jump units)
  116,497,881 Sell 300 Heating Units
  116,296,881 Buy 300 Medical Units
  116,296,881 Jump from Hothor to Sickonia (2 hyper jump units)
  118,667,481 Sell 300 Medical Units
  118,217,481 Buy 300 Heating Units
  118,217,481 Jump from Sickonia to Hothor (2 hyper jump units)
  120,647,481 Sell 300 Heating Units
  120,446,481 Buy 300 Medical Units
  120,446,481 Jump from Hothor to HomeWorld (2 hyper jump units)
  122,853,381 Sell 300 Medical Units
' --funds 100000000 --start Earth --flight_plan Volcana,Richiana,Eden,Tribonia,Medoca,Sickonia,Earth,Schooloria

t '   99,550,000 Buy 300 Heating Units
   99,550,000 Jump from Earth to Richiana (1 hyper jump units)
  101,834,200 Sell 300 Heating Units
  101,451,100 Buy 300 Novelty Packs
  101,451,100 Jump from Richiana to Tribonia (1 hyper jump units)
  102,999,700 Sell 300 Novelty Packs
  102,719,200 Buy 300 Medical Units
  102,719,200 Jump from Tribonia to Sickonia (1 hyper jump units)
  105,089,800 Sell 300 Medical Units
  104,639,800 Buy 300 Heating Units
  104,639,800 Jump from Sickonia to Hothor (2 hyper jump units)
  107,069,800 Sell 300 Heating Units
  106,868,800 Buy 300 Medical Units
  106,868,800 Jump from Hothor to Loony (1 hyper jump units)
  108,384,700 Sell 300 Medical Units
  108,294,400 Buy 300 Clothes Bundles
  108,294,400 Jump from Loony to Dreamora (1 hyper jump units)
  109,397,500 Sell 300 Clothes Bundles
  109,267,000 Buy 300 Medical Units
  109,267,000 Jump from Dreamora to Eden (1 hyper jump units)
  109,402,000 Sell 300 Medical Units
  109,352,000 Buy 2 Eden Warp Units
  109,325,776 Buy 298 Tree Growth Kits
  109,325,776 Eden warp from Eden to Zoolie
  111,428,464 Sell 298 Tree Growth Kits
  111,311,854 Buy 299 Medical Units
  111,311,854 Eden warp from Zoolie to Sickonia
  113,674,552 Sell 299 Medical Units
  113,224,552 Buy 300 Heating Units
  113,224,552 Jump from Sickonia to Hothor (2 hyper jump units)
  115,654,552 Sell 300 Heating Units
  115,453,552 Buy 300 Medical Units
  115,453,552 Jump from Hothor to Sickonia (2 hyper jump units)
  117,824,152 Sell 300 Medical Units
  117,374,152 Buy 300 Heating Units
  117,374,152 Jump from Sickonia to Hothor (2 hyper jump units)
  119,804,152 Sell 300 Heating Units
  119,603,152 Buy 300 Medical Units
  119,603,152 Jump from Hothor to HomeWorld (2 hyper jump units)
  122,010,052 Sell 300 Medical Units
' --funds 100000000 --start Earth --flight_plan Richiana,Tribonia,Sickonia,Uniland,StockWorld,Loony,Dreamora,Eden

t '   99,895,000 Buy 300 Ground Weapons
   99,895,000 Jump from Earth to Tribonia (1 hyper jump units)
  101,740,300 Sell 300 Ground Weapons
  101,557,900 Buy 300 Tree Growth Kits
  101,557,900 Jump from Tribonia to Dune (1 hyper jump units)
  103,480,300 Sell 300 Tree Growth Kits
  103,480,300 Jump from Dune to Eden (1 hyper jump units)
  103,430,300 Buy 2 Eden Warp Units
  103,404,076 Buy 298 Tree Growth Kits
  103,404,076 Jump from Eden to Medoca (1 hyper jump units)
  104,405,952 Sell 298 Tree Growth Kits
  104,371,682 Buy 298 Medical Units
  104,371,682 Jump from Medoca to HomeWorld (1 hyper jump units)
  106,762,536 Sell 298 Medical Units
  106,519,964 Buy 298 Plastic Trinkets
  106,519,964 Jump from HomeWorld to Baboria (1 hyper jump units)
  108,346,406 Sell 298 Plastic Trinkets
  108,197,406 Buy 298 Medical Units
  108,197,406 Eden warp from Baboria to HomeWorld
  110,588,260 Sell 298 Medical Units
  110,422,913 Buy 299 Clothes Bundles
  110,422,913 Jump from HomeWorld to Dreamora (1 hyper jump units)
  111,522,336 Sell 299 Clothes Bundles
  111,392,271 Buy 299 Medical Units
  111,392,271 Eden warp from Dreamora to Sickonia
  113,754,969 Sell 299 Medical Units
  113,304,969 Buy 300 Heating Units
  113,304,969 Jump from Sickonia to Hothor (1 hyper jump units)
  115,734,969 Sell 300 Heating Units
  115,533,969 Buy 300 Medical Units
  115,533,969 Jump from Hothor to Sickonia (2 hyper jump units)
  117,904,569 Sell 300 Medical Units
  117,454,569 Buy 300 Heating Units
  117,454,569 Jump from Sickonia to Hothor (2 hyper jump units)
  119,884,569 Sell 300 Heating Units
  119,683,569 Buy 300 Medical Units
  119,683,569 Jump from Hothor to Sickonia (2 hyper jump units)
  122,054,169 Sell 300 Medical Units
  121,604,169 Buy 300 Heating Units
  121,604,169 Jump from Sickonia to Hothor (2 hyper jump units)
  124,034,169 Sell 300 Heating Units
' --funds 100000000 --start Earth --flight_plan Tribonia,Dune,Eden,Medoca,HomeWorld,Baboria,Dreamora,Hothor

t '   99,550,000 Buy 300 Heating Units
   99,550,000 Jump from Earth to Hothor (2 hyper jump units)
  101,980,000 Sell 300 Heating Units
  101,779,000 Buy 300 Medical Units
  101,779,000 Jump from Hothor to HomeWorld (1 hyper jump units)
  104,185,900 Sell 300 Medical Units
  103,941,700 Buy 300 Plastic Trinkets
  103,941,700 Jump from HomeWorld to Baboria (1 hyper jump units)
  105,780,400 Sell 300 Plastic Trinkets
  105,630,400 Buy 300 Medical Units
  105,630,400 Jump from Baboria to Sickonia (2 hyper jump units)
  108,001,000 Sell 300 Medical Units
  107,551,000 Buy 300 Heating Units
  107,551,000 Jump from Sickonia to Hothor (2 hyper jump units)
  109,981,000 Sell 300 Heating Units
  109,780,000 Buy 300 Medical Units
  109,780,000 Jump from Hothor to Sickonia (2 hyper jump units)
  112,150,600 Sell 300 Medical Units
  111,700,600 Buy 300 Heating Units
  111,700,600 Jump from Sickonia to Hothor (2 hyper jump units)
  114,130,600 Sell 300 Heating Units
  113,929,600 Buy 300 Medical Units
  113,929,600 Jump from Hothor to Sickonia (2 hyper jump units)
  116,300,200 Sell 300 Medical Units
  115,850,200 Buy 300 Heating Units
  115,850,200 Jump from Sickonia to Hothor (2 hyper jump units)
  118,280,200 Sell 300 Heating Units
' --funds 100000000 --start Earth --flight_plan 'HugeLind Mar,Dreamora,HomeWorld,Baboria,Dogafetch,StockWorld,Volcana,Gojuon'

t '   99,850,000 Buy 300 Medical Units
   99,850,000 Jump from Earth to Sickonia (2 hyper jump units)
  102,220,600 Sell 300 Medical Units
  101,770,600 Buy 300 Heating Units
  101,770,600 Jump from Sickonia to Hothor (2 hyper jump units)
  104,200,600 Sell 300 Heating Units
  103,999,600 Buy 300 Medical Units
  103,999,600 Jump from Hothor to Sickonia (2 hyper jump units)
  106,370,200 Sell 300 Medical Units
  105,920,200 Buy 300 Heating Units
  105,920,200 Jump from Sickonia to Hothor (2 hyper jump units)
  108,350,200 Sell 300 Heating Units
  108,149,200 Buy 300 Medical Units
  108,149,200 Jump from Hothor to Sickonia (2 hyper jump units)
  110,519,800 Sell 300 Medical Units
  110,069,800 Buy 300 Heating Units
  110,069,800 Jump from Sickonia to Hothor (2 hyper jump units)
  112,499,800 Sell 300 Heating Units
  112,298,800 Buy 300 Medical Units
  112,298,800 Jump from Hothor to Sickonia (2 hyper jump units)
  114,669,400 Sell 300 Medical Units
  114,219,400 Buy 300 Heating Units
  114,219,400 Jump from Sickonia to Hothor (2 hyper jump units)
  116,649,400 Sell 300 Heating Units
' --funds 100000000 --start Earth --flight_plan 'StockWorld,Norhaven,Metallica,Plague,WeaponWorld,Dogafetch,Baboria,Uniland'

t '      850,000 Buy 300 Medical Units
      850,000 Jump from Earth to HomeWorld (2 hyper jump units)
    3,256,900 Sell 300 Medical Units
    2,518,900 Buy 300 Jewels
    2,518,900 Jump from HomeWorld to WeaponWorld (2 hyper jump units)
    4,068,700 Sell 300 Jewels
           28 Buy 21191 Fighter Drones

'$'\r''Drones were 245.41 each
' --funds 1000000 --start Earth --fuel 4 --drones 21191

impossible="$(mktemp)"
t '' --funds 1000000 --start Earth --fuel 4 --drones 21192 2> "$impossible"
grep -q 'Cannot acheive success criteria' "$impossible" || exit_status=1
rm "$impossible"

exit "$exit_status"
