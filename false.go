// Copyright 2019 Matt Martz <matt@sivel.net>
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
)

func main() {
	flag.Usage = func() {
		fmt.Println(`Usage: false [ignored command line arguments]
  or:  false OPTION
Exit with a status code indicating success.

      --help     display this help and exit

NOTE: your shell may have its own version of false, which usually supersedes
the version described here.  Please refer to your shell's documentation
for details about the options it supports.`)
	}
	flag.Parse()
	os.Exit(1)
}
