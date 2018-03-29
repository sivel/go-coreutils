// Copyright 2018 Matt Martz <matt@sivel.net>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

/*const (
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
)*/

var iso8601Map = map[string]string{
	//"date": "2006-01-02T15:04:05Z07:00"
	"date":    "2006-01-02",
	"hours":   "2006-01-02T15-0700",
	"minutes": "2006-01-02T15:04-0700",
	"seconds": "2006-01-02T15:04:05-0700",
	"ns":      "2006-01-02T15:04:05.000000000-0700",
}

var rfc3339Map = map[string]string{
	"date":    "2006-01-02",
	"seconds": "2006-01-02 15:04:05-0700",
	"ns":      "2006-01-02 15:04:05.000000000-0700",
}

var strftimeMap = map[string]string{
	"B":  "January",
	"b":  "Jan",
	"-m": "1",
	"m":  "01",
	"A":  "Monday",
	"a":  "Mon",
	"-d": "2",
	//stdUnderDay                                    // "_2"
	"d":  "02",
	"-H": "15",
	"H":  "15",
	"I":  "03",
	"-I": "3",
	"-M": "4",
	"M":  "04",
	"-S": "5",
	"S":  "05",
	"Y":  "2006",
	"y":  "06",
	"p":  "PM",
	//stdpm                                          // "pm"
	"Z": "MST",
	//stdISO8601TZ                                   // "Z0700"  // prints Z for UTC
	//stdISO8601SecondsTZ                            // "Z070000"
	//stdISO8601ShortTZ                              // "Z07"
	//stdISO8601ColonTZ                              // "Z07:00" // prints Z for UTC
	//stdISO8601ColonSecondsTZ                       // "Z07:00:00"
	"z": "-0700", // always numeric
	//stdNumSecondsTz                                // "-070000"
	//stdNumShortTZ                                  // "-07"    // always numeric
	//stdNumColonTZ                                  // "-07:00" // always numeric
	//stdNumColonSecondsTZ                           // "-07:00:00"
	"f": ".000000",
	//stdFracSecond9                                 // ".9", ".99", ..., trailing zeros omitted

	"c": time.UnixDate,
}

var strftimeRe = regexp.MustCompile("([^%]+)?%([a-zA-Z-]{1,2})([^%]+)?")

func convertStrftime(format string) string {
	var out []string
	for _, match := range strftimeRe.FindAllStringSubmatch(format, -1) {
		out = append(out, match[1])
		if val, ok := strftimeMap[match[2]]; ok {
			out = append(out, val)
		} else {
			out = append(out, match[2])
		}
		out = append(out, match[3])
	}
	return strings.Join(out, "")
}

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

func (sf *stringFlag) IsBoolFlag() bool {
	return true
}

func main() {
	var theTime time.Time
	var format string = "Mon Jan _2 15:04:05 MST 2006"
	var iso8601 stringFlag

	utc := flag.Bool("u", false, "print Coordinated Universal Time (UTC)")
	inDate := flag.String("date", "", "display time described by STRING, not 'now'")
	flag.Var(
		&iso8601,
		"I",
		"output date/time in ISO 8601 format. FMT='date' for date only (the default), 'hours', 'minutes', 'seconds', or 'ns' for date and time to the indicated precision. Example: 2006-08-14T02:34:56-0600",
	)
	rfc3339 := flag.String("rfc-3339", "", "")
	rfc2822 := flag.Bool("rfc-2822", false, "")
	flag.Parse()

	if iso8601.set {
		if iso8601.value == "true" {
			iso8601.Set("date")
		}
		format = iso8601Map[iso8601.String()]
	}

	if *rfc3339 != "" {
		format = rfc3339Map[*rfc3339]
	}

	if *rfc2822 {
		format = "Mon, 02 Jan 2006 15:04:05 -0700"
	}

	args := flag.Args()

	if len(args) > 1 {
		fmt.Printf("date: extra operand ‘%s’\n", args[1])
		os.Exit(1)
	} else if !iso8601.set && *rfc3339 == "" && !*rfc2822 && len(args) == 1 {
		format = convertStrftime(args[0])
	}

	if *inDate != "" {
		fmt.Println("FIXME")
	} else {
		theTime = time.Now()
	}

	if *utc {
		fmt.Println(theTime.UTC().Format(format))
	} else {
		fmt.Println(theTime.Format(format))
	}
}
