# Vochain Block Explorer 
### Deps
Golang 1.14+
## Running
#### 1. Frontend 
On the commandline navigate into the directory `./frontend` and run `go generate`, that will create `main.wasm` file in `./static`.

#### 2. Static files
Make a copy of the `wasm_exec.js` file from `$GOROOT/misc/wasm/` directory and put it in the `./static` directory.  This must be from the same golang version that you used to build `main.wasm`.

#### 3. Running
After steps 1 and 2, on the commandline navigate into the directory containing `main.go` and run `go run main.go`. Then in your favorite web browser navigate to `localhost:8081`.

----
Using [vectyUI](https://github.com/nathanhack/vectyUI) and [marwanio](https://github.com/marwan-at-work/marwanio) as inspiration