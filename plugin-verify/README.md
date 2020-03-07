# Validation of Plugins

## Build Plugin
```
# version is pre determined
go build -buildmode=plugin -o ./testdata/plugins/someplugin-r13_2-$(go version | cut -d" " -f3)-$(go env GOARCH).so someplugin/main.go
```

## Generate Signature of Plugin
```
# the filename here depends on your environment and is produced from above
go run signer/signer.go -input ./testdata/plugins/someplugin-r13_2-go1.13.4-amd64.so

### output
writing signature file: ./testdata/plugins/someplugin-r13_2-go1.13.4-amd64.so.sig
```

## Run the Wrapper
```
go run wrapper.go 

### output
...
using plugin './testdata/plugins/someplugin-r13_2-go1.13.4-amd64.so' sig './testdata/plugins/someplugin-r13_2-go1.13.4-amd64.so.sig'
plugin signature has been validated!
plugin version (R13_2) validated!
Doing Something...

```

## (Optional) Serve the Files over HTTP!
```
# run an http server to serve the files produced from above
docker run -d --rm --name plugin-server -p 18675:80 -v $(pwd)/testdata/plugins:/usr/share/nginx/html/plugins:ro nginx

# now run the wrapper again
go run wrapper.go 

### output
fetching plugin via URL: http://127.0.0.1:18675/plugins/someplugin-r13_2-go1.13.4-amd64.so
fetching plugin signature via URL: http://127.0.0.1:18675/plugins/someplugin-r13_2-go1.13.4-amd64.so.sig
plugin signature has been validated!
plugin version (R13_2) validated!
Doing Something...
```