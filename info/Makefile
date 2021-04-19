generate: process.png

process.png: process.puml 
	cat $< | docker run --rm -i think/plantuml -tpng > $@_
	mv $@_ $@
