package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path"
	"plugin"
	"runtime"
	"strings"
)

const (
	// key used to validate the authenticity of the plugin
	pubKeyPem = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEm2R44sjCU5RZzkBnpCaFXakB6iBh
0mqennUQBJ0g8BU7M1nxbecKQ+hL+kF2kBxal+/fdgeOLf5W/kCkQ3O0mw==
-----END PUBLIC KEY-----
`
	// version string used to verify that the plugin version matches what we expected
	version = "R13_2"
)

var (
	pluginName    = "someplugin"
	pluginBaseDir = "./testdata/plugins"
	// more specific paths for production
	preferredPluginFilename      = fmt.Sprintf("%s-%s-%s-%s.so", pluginName, strings.ToLower(version), runtime.Version(), runtime.GOARCH)
	preferredPluginPath          = fmt.Sprintf("%s/%s", pluginBaseDir, preferredPluginFilename)
	preferredPluginSignaturePath = fmt.Sprintf("%s.sig", preferredPluginPath)
	// maybe less specific paths for dev??
	fallbackPluginPath          = fmt.Sprintf("%s/%s.so", pluginBaseDir, pluginName)
	fallbackPluginSignaturePath = fmt.Sprintf("%s.sig", fallbackPluginPath)

	// baseurl to fetch plugin
	baseUrl = "http://127.0.0.1:18675/plugins"
	// this is where we store our downloads
	downloadDir = "./testdata/plugins-fetched"
)

type SomePlugin interface {
	Version() string
	DoSomething()
}

func downloadFile(url string) (file string, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	httpClient := http.Client{}

	// try to download the file
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Create the file
	downloadPath := fmt.Sprintf("%s/%s", downloadDir, path.Base(req.URL.Path))
	out, err := os.Create(downloadPath)
	if err != nil {
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return downloadPath, err
}

func fetchPlugin() (plugin, pluginSig string, err error) {
	// download the plugin
	pluginUrlPath := fmt.Sprintf("%s/%s", baseUrl, preferredPluginFilename)
	fmt.Println("fetching plugin via URL:", pluginUrlPath)
	pluginPath, err := downloadFile(pluginUrlPath)
	if err != nil {
		return
	}

	// download the signature file
	fmt.Println("fetching plugin signature via URL:", pluginUrlPath+".sig")
	pluginSigPath, err := downloadFile(pluginUrlPath + ".sig")
	if err != nil {
		return
	}

	return pluginPath, pluginSigPath, nil
}

func findPlugin() (plugin, pluginSig string) {
	// check for expected files
	_, errPP := os.Stat(preferredPluginPath)
	_, errPS := os.Stat(preferredPluginSignaturePath)
	switch {
	case errPP == nil && errPS == nil:
		plugin = preferredPluginPath
		pluginSig = preferredPluginSignaturePath
		fmt.Printf("using plugin '%s' sig '%s'\n", plugin, pluginSig)
		return
	case errPP == nil && errPS != nil:
		panic(fmt.Sprintf("could not find signature file '%s' (%s)\n", preferredPluginSignaturePath, errPS))
	}

	// fallback logic (dev mode??)
	_, errFP := os.Stat(fallbackPluginPath)
	_, errFS := os.Stat(fallbackPluginSignaturePath)
	switch {
	case errFP == nil && errFS == nil:
		plugin = fallbackPluginPath
		pluginSig = fallbackPluginSignaturePath
		fmt.Printf("using fallback plugin '%s' sig '%s'\n", plugin, pluginSig)
	case errFP != nil:
		panic(fmt.Sprintf("could not find plugin '%s' (%s)\n", fallbackPluginPath, errFP))
	case errFS != nil:
		panic(fmt.Sprintf("could not find signature file '%s' (%s)\n", fallbackPluginSignaturePath, errFS))
	}
	return
}

func readSignature(sigFile string) (signature string) {
	sigFileBytes, err := ioutil.ReadFile(sigFile)
	if err != nil {
		panic(err)
	}
	return string(sigFileBytes)
}

func loadPublicKey(publicKey string) (*ecdsa.PublicKey, error) {
	// decode the key, assuming it's in PEM format
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("Failed to decode PEM public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Failed to parse ECDSA public key")
	}
	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	}
	return nil, errors.New("Unsupported public key type")
}

func validateSignature(signature string, data []byte) bool {
	pubKey, err := loadPublicKey(pubKeyPem)
	if err != nil {
		fmt.Println("loadPublicKey err:", err)
		return false
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		fmt.Println("hex.DecodeString err:", err)
		return false
	}

	var esig struct {
		R, S *big.Int
	}
	if _, err := asn1.Unmarshal(signatureBytes, &esig); err != nil {
		fmt.Println("asn1.Unmarshal", err)
		return false
	}
	digest := sha256.Sum256(data)
	return ecdsa.Verify(pubKey, digest[:], esig.R, esig.S)
}

func main() {
	// attempt to fetch from URL first
	pluginFile, sigFile, err := fetchPlugin()
	if err != nil {
		fmt.Println("unable to download plugin:", err)
		// fallback to local files
		fmt.Println("falling back to loading plugins from local filesystem")
		pluginFile, sigFile = findPlugin()
	}

	// validate the plugin
	pluginFileBytes, err := ioutil.ReadFile(pluginFile)
	if err != nil {
		panic(err)
	}
	validate := validateSignature(readSignature(sigFile), pluginFileBytes)
	if !validate {
		panic("plugin signature is not valid")
	}
	fmt.Println("plugin signature has been validated!")

	// the module signature is OK, should be safe to load now
	userccPlugin, err := plugin.Open(pluginFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// some compatibility test checks...
	symUsercc, err := userccPlugin.Lookup("SomePlugin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var usercc SomePlugin
	usercc, ok := symUsercc.(SomePlugin)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	// validate the version
	if strings.ToLower(usercc.Version()) != strings.ToLower(version) {
		fmt.Printf("version validatin failure: expected %s but got %s\n", version, usercc.Version())
		os.Exit(1)
	}
	fmt.Printf("plugin version (%s) validated!\n", version)

	// finally continue to do something with this plugin...
	usercc.DoSomething()
}
