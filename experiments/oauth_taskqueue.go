// Must be run specificaly, `go run oauth.go`
// Scope can be found here via the discovery service:
//   https://www.googleapis.com/discovery/v1/apis/taskqueue/v1beta2/rest

// ```bash
// go run oauth.go \
// -s ../vault/voom-builder.json
// -k ../vault/voom-builder-4d7e7eb89abc.pem
//
// ```

package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strings"

    "code.google.com/p/goauth2/oauth/jwt"
)

var scopes = []string{
    "https://www.googleapis.com/auth/taskqueue",
    "https://www.googleapis.com/auth/taskqueue.consumer",
}

var scope = strings.Join(scopes, " ")

var (
    secretsFile = flag.String("s", "", "JSON encoded secrets for the service account")
    // pemFile     = flag.String("k", "", "private pem key file for the service account")
)

const usageMsg = `
You must specify -k and -s.

To obtain client secrets and pem, see the "OAuth 2 Credentials" section under
the "API Access" tab on this page: https://code.google.com/apis/console/

Google Cloud Storage must also be turned on in the API console.
`

func main() {
    flag.Parse()

    if *secretsFile == "" { //|| *pemFile == "" {
        flag.Usage()
        fmt.Println(usageMsg)
        return
    }

    // Read the secret file bytes into the config.
    secretBytes, err := ioutil.ReadFile(*secretsFile)
    if err != nil {
        log.Fatal("error reading secerets file:", err)
    }
    var config struct {
        ClientEmail string `json:"client_email"`
        ClientID    string `json:"client_id"`
        TokenURI    string `json:"token_uri"`
        PrivateKey  string `json:"private_key"`
    }
    err = json.Unmarshal(secretBytes, &config)
    if err != nil {
        log.Fatal("error unmarshalling secerets:", err)
    }

    // Get the project ID from the client ID.
    projectID := strings.SplitN(config.ClientID, "-", 2)[0]
    // projectID := "voom-registry-service"

    // Read the pem file bytes for the private key.
    // keyBytes, err := ioutil.ReadFile(*pemFile)
    // if err != nil {
    //     log.Fatal("error reading private key file:", err)
    // }
    keyBytes := []byte(config.PrivateKey)

    // Craft the ClaimSet and JWT token.
    t := jwt.NewToken(config.ClientEmail, scope, keyBytes)
    t.ClaimSet.Aud = config.TokenURI

    // We need to provide a client.
    c := &http.Client{}

    // Get the access token.
    o, err := t.Assert(c)
    if err != nil {
        log.Fatal("assertion error:", err)
    }

    // Refresh token will be missing, but this access_token will be good
    // for one hour.
    fmt.Printf("access_token = %v\n", o.AccessToken)
    fmt.Printf("refresh_token = %v\n", o.RefreshToken)
    fmt.Printf("expires %v\n", o.Expiry)

    // voom-registry-service
    url := "https://www.googleapis.com/taskqueue/v1beta2/projects/voom-registry-service/taskqueues/jobs-pending/tasks"
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal("http.NewRequest:", err)
    }
    req.Header.Set("Authorization", "OAuth "+o.AccessToken)
    req.Header.Set("x-goog-api-version", "2")
    req.Header.Set("x-goog-project-id", projectID)

    // Make the request.
    r, err := c.Do(req)
    if err != nil {
        log.Fatal("API request error:", err)
    }
    defer r.Body.Close()

    // Write the response to standard output.
    res, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal("error reading API request results:", err)
    }
    fmt.Printf("\nRESULT:\n%s\n", res)
}
