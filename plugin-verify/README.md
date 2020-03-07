# Validation of Plugins

## Build Plugin
```
go build -buildmode=plugin -o ./testdata/plugins/someplugin.so someplugin/main.go
```

## Generate Signature of Plugin
```
go run signer/signer.go -input ./testdata/plugins/someplugin.so

### output
writing signature file: ./testdata/plugins/someplugin.so.sig
```

## Run the Wrapper
```
go run wrapper.go 

### output
using fallback plugin './testdata/plugins/someplugin.so' sig './testdata/plugins/someplugin.so.sig'
plugin signature has been validated!
plugin version (R13_2) validated!
Doing Something...
```