module github.com/newrelic/nri-mongodb

go 1.18

require (
	github.com/globalsign/mgo v0.0.0-20190517090918-73267e130ca1
	github.com/newrelic/infra-integrations-sdk v3.7.3+incompatible
	github.com/stretchr/testify v1.7.5
	github.com/xeipuuv/gojsonschema v1.2.0
)

replace github.com/globalsign/mgo => github.com/mhill-anynines/mgo v0.0.0-20190116155901-0d5878c1bace

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
