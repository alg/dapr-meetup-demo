dapr run \
	--app-id warehouse \
	--app-port 8001 \
	--dapr-http-port 3501 \
    --resources-path ../components \
    --resources-path ./components \
	-- \
	go run .
