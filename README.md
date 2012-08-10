Poker
======

This is a library for creating agents that can play Texas Hold'em and includes
an example agent Hob2 who can play randomly or to maximize EV based on
seven-card hand strength. Currently, the library implements only the [ACPC
protocol][1], but support for more should not be hard to add. There is also
code for generating game trees and calculating a Nash Equilibrium, but it is
not complete.

[1]: http://www.computerpokercompetition.org/index.php?option=com_rokdownloads&view=file&task=download&id=130:acpc-2011-protocol


Dependencies
------------

The library uses the [2 + 2 hand ranks table][2] to evaluate hands. The table
table can be generated both in [Windows][3] or [Linux][4] using Wine. Hob2
expects the table to be in the same directory it is run in.

[2]: http://archives1.twoplustwo.com/showflat.php?Cat=0&Number=8513906&amp;amp;amp;page=2&fpart=1&vc=1
[3]: http://www.codingthewheel.com/archives/poker-hand-evaluator-roundup#2p2
[4]: https://github.com/davekong/two-plus-two-table-Linux

License
--------

    Poker poker library
    Copyright (C) 2012 David Campbell <dcampbell24@gmail.com>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
