package userstatsapi

//go:generate protoc -I ../../../api --go_out=plugins=grpc:. UserStats.proto
