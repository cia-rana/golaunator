module github.com/cia-rana/golaunator/example

require (
	github.com/awalterschulze/gographviz v0.0.0-20181013152038-b2885df04310
	github.com/cia-rana/golaunator/fortune v0.0.0
	github.com/cia-rana/golaunator/iia v0.0.0
	github.com/fogleman/gg v1.1.0
)

replace github.com/cia-rana/golaunator/iia => ../iia

replace github.com/cia-rana/golaunator/fortune => ../fortune
