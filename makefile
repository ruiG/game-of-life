.PHONY: run fmt package clean

run: # runs the program
	go run main.go

fmt: # format Go code
	go fmt ./...

package: # bundle the app for macOS
	go run fyne.io/tools/cmd/fyne@latest package -os darwin -icon gol.png -release

clean: # remove build artifacts
	rm -rf game-of-life.app