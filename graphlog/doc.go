/*
graphlog is a tool for visualizing the log of a poker game.

graphlog uses gnuplot[1] to generate an interactive graph displaying the scores of
each player over the course of the log. graphlog assumes that the entries in the
log represent consective hands and numbers them automatically.

[1]: http://www.gnuplot.info/

Usage:

	graphlog log

Gnuplot graph window commands:

	2x<B1>             print coordinates to clipboard using `clipboardformat`
	                   (see keys '3', '4')
	<B2>               annotate the graph using `mouseformat` (see keys '1', '2')
	                   or draw labels if `set mouse labels is on`
	<Ctrl-B2>          remove label close to pointer if `set mouse labels` is on
	<B3>               mark zoom region (only for 2d-plots and maps).
	<B1-Motion>        change view (rotation). Use <ctrl> to rotate the axes only.
	<B2-Motion>        change view (scaling). Use <ctrl> to scale the axes only.
	<Shift-B2-Motion>  vertical motion -- change xyplane
	<wheel-up>         scroll up (in +Y direction).
	<wheel-down>       scroll down.
	<shift-wheel-up>   scroll left (in -X direction).
	<shift-wheel-down>  scroll right.
	<control-wheel-up>  zoom in toward the center of the plot.
	<control-wheel-down>   zoom out.
	<shift-control-wheel-up>  zoom in only the X axis.
	<shift-control-wheel-down>  zoom out only the X axis.

	Space          raise gnuplot console window
	q            * close this plot window

	a              `builtin-autoscale` (set autoscale keepfix; replot)
	b              `builtin-toggle-border`
	e              `builtin-replot`
	g              `builtin-toggle-grid`
	h              `builtin-help`
	l              `builtin-toggle-log` y logscale for plots, z and cb for splots
	L              `builtin-nearest-log` toggle logscale of axis nearest cursor
	m              `builtin-toggle-mouse`
	r              `builtin-toggle-ruler`
	1              `builtin-previous-mouse-format`
	2              `builtin-next-mouse-format`
	3              `builtin-decrement-clipboardmode`
	4              `builtin-increment-clipboardmode`
	5              `builtin-toggle-polardistance`
	6              `builtin-toggle-verbose`
	7              `builtin-toggle-ratio`
	n              `builtin-zoom-next` go to next zoom in the zoom stack
	p              `builtin-zoom-previous` go to previous zoom in the zoom stack
	u              `builtin-unzoom`
	Right          `builtin-rotate-right` only for splots; <shift> increases amount
	Up             `builtin-rotate-up` only for splots; <shift> increases amount
	Left           `builtin-rotate-left` only for splots; <shift> increases amount
	Down           `builtin-rotate-down` only for splots; <shift> increases amount
	Escape         `builtin-cancel-zoom` cancel zoom region

	             * indicates this key is active from all plot windows
*/
package documentation
