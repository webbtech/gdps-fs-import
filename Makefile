include .env

default: build

clean:
	@rm -rf dist
	@mkdir -p dist

build: clean
	@for dir in `ls handler`; do \
		GOOS=linux go build -o dist/$$dir github.com/pulpfree/gales-fuelsale-export/handler/$$dir; \
	done
	cp ./config/defaults.yaml dist/

run:
	sam local start-api -n env.json