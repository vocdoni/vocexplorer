# Vochain Block Explorer

Deployed at [https://explorer.vocdoni.net/](https://explorer.vocdoni.net/) and [https://explorer.dev.vocdoni.net/](https://explorer.dev.vocdoni.net/)

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

On the commandline navigate into the directory `./frontend` and run 
~~~
env GOARCH=wasm GOOS=js go build -ldflags "-s -w" -trimpath -o ../static/main.wasm
~~~
to compile the frontend code into wasm.

#### 2. Static files

Make a copy of the `wasm_exec.js` file from `$GOROOT/misc/wasm/` directory and put it in the `./static` directory.  This must be from the same golang version that you used to build `main.wasm`.

Get back to the root path and run `yarn` to install and compile the required style assets.

If you want to renew the styles, run `yarn gulp`, or in case you wanna watch for changes, `yarn gulp sass:watch`.

There's also a gulp task for watching `.go` files changes in `./frontend` files and regenerating the `main.wasm` file: `yarn gulp go:watch`.

You can watch both `.go` and `.scss` file changes by just using

~~~bash
yarn gulp watch
~~~

#### 4. Backend

After steps 1, and 2, navigate back into `vocexplorer` and run
~~~ 
go run main.go
~~~ 
Then in your favorite web browser navigate to localhost at the specified port.

Options for `main.go`:
- `--dataDir` `(string)`             directory where data is stored (default "/Users/natewilliams/.vocexplorer")
- `--refreshTime` `(int)`            number of seconds between each content refresh (default 10)
- `--gatewayUrl` `(string)`          vocdoni node URL to query for data
- `--disableGzip`                    use to disable gzip compression on web server
- `--hostURL` `(string)`             url to host block explorer (default "http://localhost:8081")
- `--logLevel` `(string)`            log level <debug, info, warn, error> (default "error")
----
