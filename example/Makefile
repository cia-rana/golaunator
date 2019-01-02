%_graph:
	go run $*_graph.go util.go | neato -n2 -T png -o img/$*.png

fortune_movie:
	go run fortune_graph.go util.go

clean:
	rm *.png
	rm ~/Projects/SimpleImageServer/*.png
	rm img/*
