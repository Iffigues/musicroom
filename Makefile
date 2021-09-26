all: build

init:
	make init -C server

build:
	@make -C server

clean:
	make clean -C server

fclean:
	make fclean -C server

fmt:
	@make fmt -C server

push:
	make clean -C server
	git add .
	git commit -m "push"
	git push

fpush:
	make fclean -C server
	git add .
	git commit -m "push"
	git push

