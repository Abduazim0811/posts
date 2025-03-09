DB_URL=postgres://postgres:Abdu0811@172.17.0.2:5432/postservice?sslmode=disable

run:
	go run cmd/main.go

migrate_up:
	migrate -path ./migrations -database ${DB_URL} -verbose up

migrate_down:
	migrate -path ./migrations -database ${DB_URL} -verbose down

migrate_force:
	migrate -path ./migrations -database ${DB_URL} -verbose force 1

migrate_file:
	migrate create -ext sql -dir migrations -seq gin

migrate_proto:
	protoc --go_out=protos/postProto --go-grpc_out=protos/postProto protos/postProto/post.proto
