# Implementation of a messenger using blockchain technology

## The application is available at http://ec2-35-160-233-232.us-west-2.compute.amazonaws.com:8080/

## Application components

- blockchain structure
- algorithms related to Blockchain technology:
    - block validation,
    - generating hash function values,
    - conflict handling
- transaction handling implementation
- user interface - CLI or simple web page

## Used libraries:
- "crypto/sha256"
- "encoding/hex"
- "fmt"
- "html/template"
- "log"
- "strconv"
- "time"
- "io"
- "net/http"
- "net/http/httptest"
- "strings"
- "testing"
- "reflect"

## To run application:
1. Go to main project package
2. run command: go run main.go
3. open browser and paste url: localhost:8080

## To run tests:
1. Go to main project package
2. run command: go test

## To run tests and generate coverage with html
1. Go to main project package
2. run command: go test -v -cover -coverprofile=results
3. run command: go tool cover -html=results
4. Browser should open with coverage test results