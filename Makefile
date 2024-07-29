run:
	@ printf "Starting Aplication... \n"
	@ go run cmd/api/main.go 
	
tidy:
	go mod tidy