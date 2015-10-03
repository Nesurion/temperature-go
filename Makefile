PROJECT=temperature-go

$(PROJECT): 
	go build

install:
	cp $(PROJECT) /usr/local/bin/$(PROJECT)

clean:
	rm $(PROJECT)
