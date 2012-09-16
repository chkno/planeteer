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

exit "$exit_status"
