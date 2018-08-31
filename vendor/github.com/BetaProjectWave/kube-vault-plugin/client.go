package vault

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	// Info logger
	Info    *log.Logger

	// Error logger
	Error   *log.Logger
)

const (
	// AUTH used in the vault REST path for login endpoint
	AUTH = "auth"

	// ClientToken the json name field of the client token we get in the login response
	ClientToken = "client_token"

	// JWT vault token header for querying the secrets
	JWT = "X-Vault-Token"

	// DATA the json field object name representing the secrets
	DATA = "data"

	// KUBERNETES used in the REST path for login request
	KUBERNETES = "kubernetes"

	// APIVersion is the vault REST api version
	APIVersion = "v1"

	// LOGIN used in the vault REST path for the login request.
	LOGIN = "login"

	// ContentType used in all Get/Post Vault REST request.
	ContentType = "application/json"
)

func init() {
	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Llongfile)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Llongfile)
}

// Config for vault
type Config struct {
	// The path to root CA used by Vault
	CaFile string

	// Certificate file Public key
	CertFile string

	// Private Key
	KeyFile string

	// The JWT vault token
	Jwt string

	// The role attached to the JWT vault token
	Role string

	// Secret path
	SecretPath string

	// Address of the Vault server, exep: https://vault.esxample.com
	Address string

	// Vault plugin enabled
	Enabled bool
}

// ClientVault used for Vault HTTP client
type ClientVault struct {
	HTTPClient *http.Client
	Data       map[string]interface{}
	Config     Config
	Token      string
}

// LoginRequest used for Vault login
type LoginRequest struct {
	Role string `json:"role"`
	Jwt  string `json:"jwt, string"`
}

// ReadSecret form the vault repository for given key
func (client *ClientVault) ReadSecret(key string) string {
	Info.Printf("Reading secret for a given key: %s", key)
	// Read secret for a given key
	return client.Data[key].(string)
}

// Perform login
func (client *ClientVault) login() (string, error) {
	jwt, err := readJWT(client.Config)
	if err != nil {
		Error.Printf("Failed to unmarshal the login request")
		return "", err
	}
	request := LoginRequest{Role: client.Config.Role, Jwt: jwt}
	body, err := json.Marshal(request)
	if err != nil {
		Error.Printf("Failed to unmarshal the login request")
		Error.Fatal(err)
	}
	resp, err := client.HTTPClient.Post(loginEndpoint(client.Config), ContentType, bytes.NewReader(body))
	if err != nil {
		Error.Printf("Failed to login to the Vault server")
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to login to the Vault server")
	}

	defer resp.Body.Close()

	// read all the bytes
	response, err := ioutil.ReadAll(resp.Body)

	byt := []byte(string(response))
	var result map[string]interface{}
	if err := json.Unmarshal(byt, &result); err != nil {
		Error.Printf("Failed to unmarshall the login response")
		return "", err
	}
	auth := result[AUTH].(map[string]interface{})
	token := auth[ClientToken].(string)
	return token, nil
}

// Load all the secrets into a  Global map
func (client *ClientVault) loadAllSecrets() (map[string]interface{}, error) {
	path := secretEndpoint(client.Config)
	Info.Printf("Loading secrets from vault server on path: %s", path)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	req.Header.Set(JWT, client.Token)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		Error.Printf("Failed to load all the secrets from the vault server: %s", err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		Error.Printf("failed to load all the secrets")
		return nil, errors.New("failed to load all the secrets")
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)

	byt := []byte(string(response))

	var result map[string]interface{}

	if err := json.Unmarshal(byt, &result); err != nil {
		Error.Printf("Failed  to unmarshal the login resposne: %s", err.Error())
		return nil, err
	}
	// initialise a global Map
	data, ok := result[DATA]
	if !ok {
		return nil, errors.New("No secret found in the given path: " + path)
	}
	return data.(map[string]interface{}), nil
}

// NewClient returns an instance of VaultClient
func NewClient(config Config) (*ClientVault, error) {
	tlsConfig, err := addTLS(config)
	if err != nil {
		Error.Println(err.Error())
		return nil, err
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	// Http client instance
	httpClient := &http.Client{Transport: transport}

	client := &ClientVault{HTTPClient: httpClient, Config: config}

	// perform login
	Info.Printf("Performing login to the vault server")
	token, err := client.login()
	if err != nil {
		Error.Printf("Failed to login to the vault server: %s", err.Error())
		return nil, err
	}
	client.Token = token
	Info.Printf("Successfully loged in to the vault server")

	// Load all the secrets
	client.Data, err = client.loadAllSecrets()
	if err != nil {
		Error.Printf("Failed to load the secrets: %s", err.Error())
		return nil, err
	}
	return client, nil
}

// Adds TLS setting for secure connection
func addTLS(config Config) (*tls.Config, error) {

	var cert tls.Certificate
	var err error

	// Load client cert
	if usePrivateKey(config) {
		cert, err = tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, err
		}
	}
	caCert, err := ioutil.ReadFile(config.CaFile)
	if err != nil {
		Error.Printf("Failed to load CA certificate for given path: %s", config.CaFile)
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	return tlsConfig, nil
}

// Helper function
func usePrivateKey(config Config) bool {
	return config.KeyFile != ""
}

// helper for secret endpoint
func secretEndpoint(config Config) string {
	var endPoint bytes.Buffer
	endPoint.WriteString(config.Address)
	endPoint.WriteString("/")
	endPoint.WriteString(APIVersion)
	endPoint.WriteString("/")
	endPoint.WriteString(config.SecretPath)
	return endPoint.String()
}

// helper for secret endpoint
func loginEndpoint(config Config) string {
	var endPoint bytes.Buffer
	endPoint.WriteString(config.Address)
	endPoint.WriteString("/")
	endPoint.WriteString(APIVersion)
	endPoint.WriteString("/")
	endPoint.WriteString(AUTH)
	endPoint.WriteString("/")
	endPoint.WriteString(KUBERNETES)
	endPoint.WriteString("/")
	endPoint.WriteString(LOGIN)
	return endPoint.String()
}

// Load the token from a given path
func readJWT(config Config) (string, error) {
	token, err := ioutil.ReadFile(config.Jwt)
	if err != nil {
		Error.Printf("Failed to load JWT token: %s", err.Error())
		return "", err
	}
	return strings.Trim(string(token), "\n"), nil
}
