dapr run \
	--app-id store \
	--app-port 8000 \
	--dapr-http-port 3500 \
    --resources-path ../components \
	-- \
	go run .
