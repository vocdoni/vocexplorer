# Vochain Block Explorer

## Deps

Golang 1.14+

## Running

### Docker

Navigate to `vocexplorer/docker/vocexplorer`
~~~
docker-compose build
docker-compose up
~~~
That's it!

### Manual

#### 1. Frontend

On the commandline navigate into the directory `./frontend` and run `go generate`, that will create `main.wasm` file in `./static`.

#### 2. Static files

Make a copy of the `wasm_exec.js` file from `$GOROOT/misc/wasm/` directory and put it in the `./static` directory.  This must be from the same golang version that you used to build `main.wasm`.

Get back to the root path and run `yarn` to install and compile the required style assets.

If you want to renew the styles, run `yarn gulp`, or in case you wanna watch for changes, `yarn gulp sass:watch`.

There's also a gulp task for watching `.go` files changes in `./frontend` files and regenerating the `main.wasm` file: `yarn gulp go:watch`.

You can watch both `.go` and `.scss` file changes by just using

~~~bash
yarn gulp watch
~~~

#### 3. Dvotenode

Clone the `go-dvote` [repository](https://gitlab.com/vocdoni/go-dvote).  
Then run 
~~~
cd go-dvote
go run cmd/dvotenode/dvotenode.go --w3Enabled=False --vochainNoWaitSync 
~~~
Make sure that the no-sync options are enabled.

#### 4. Vocexplorer

After steps 1, 2, and 3, navigate back into `vocexplorer` and run
~~~ 
go run main.go
~~~ 
Then in your favorite web browser navigate to localhost at the specified port.

Options for `main.go`:
- `--disableGzip`          use to disable gzip compression on web server
- `--gatewayHost` `(string)`   gateway API host to connect to (default "ws://0.0.0.0:9090/dvote")
- `--hostURL` `(string)`       url to host block explorer (default "http://localhost:8081")
- `--logLevel` `(string)`      log level <debug, info, warn, error> (default "error")
- `--refreshTime` `(int)`      Number of seconds between each content refresh (default 5)
- `--vochainHost` `(string)`   gateway API host to connect to (default "http://0.0.0.0:26657")

----
