dapr run \
	--app-id notifications \
	--app-port 8002 \
	--app-protocol grpc \
	--dapr-http-port 3502 \
    --resources-path ../components \
    --resources-path ./components \
	-- \
	go run .
