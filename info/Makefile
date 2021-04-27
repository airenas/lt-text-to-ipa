generate: process.png processGET.png

%.png: %.puml 
	cat $< | docker run --rm -i think/plantuml -tpng > $@_
	mv $@_ $@
