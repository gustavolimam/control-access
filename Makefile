install:
	cd client && yarn install 
	
build:
	cd client && yarn build

run: 
	go run src/main.go