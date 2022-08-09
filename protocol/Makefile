.PHONY : all health

all: health

health:
	protoc -I . health.proto --go_out=plugins=grpc:../protogo/health/

member:
	protoc -I . member.proto --go_out=plugins=grpc:../protogo/member/
