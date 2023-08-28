package uxt

//import (
//"bytes"
//"encoding/json"
//"io"
//"io/ioutil"
//"net/http"
//"strconv"

//"github.com/gocql/gocql"
//"gitlab.com/textfridayy/uno/uf"
//)

//type UfRequest struct {
//Url  string
//Body map[string]interface{}
//}

//type UfResponse struct {
//Resp    *http.Response
//Body    map[string]interface{}
//BodyObj interface{}
//UfId    gocql.UUID
//Errors  []uf.ApiError
//}

//type AbstractBody struct {
//BodyObj interface{}
//Errors  []uf.ApiError
//}

//type ApiErrorsFromBody struct {
//Errors []uf.ApiError
//}

//func MakeRequest(
//gibs *uf.Gibs,
//address string,
//path string,
//pkg interface{},
//out interface{},
//) (ures UfResponse) {
//url := MakeUrl(address, path, gibs.Conf)
//pkgString, err := json.Marshal(pkg)
//if err != nil {
//ures.AddError(uf.EncodeUfRequestBodyError)
//}

//uf.Debug("=====================================")
//uf.Debug("Internal request to " + url.String() + " with")
//uf.Debug(pkg)

//client := &http.Client{}
//req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(pkgString))
//req.Header.Set("UfAuth", gibs.UfAuth)
//req.Header.Set("Content-Type", "application/json")
//req.Header.Set("UfId", gibs.UfId.String())
//req.Header.Set("UfKey", gibs.Conf.UfKey)

//ures.Resp, err = client.Do(req)
//if err != nil {
//ures.AddError(uf.CustomError(err))
//}

//data, err := ioutil.ReadAll(ures.Resp.Body)
//if err != nil {
//panic(err.Error())
//}

//if ures.Resp == nil {
//ures.Body = make(map[string]interface{})
//ures.Errors = []uf.ApiError{}
//} else {
//uf.PrintResponseBody(data)

//reader := bytes.NewReader(data)
//err = json.NewDecoder(reader).Decode(&out)
//if err != nil {
//ures.AddError(uf.DecodeUfResponseBodyError)
//}

//reader = bytes.NewReader(data)
//ures.unpackUfResponseErrors(reader)
//ures.BodyObj = out
//}

//uf.Debug(ures.BodyObj)
//ures.LogErrors()

//return ures
//}

//func (ures *UfResponse) LogErrors() {
//if ures.Errors != nil {
//for _, ufError := range ures.Errors {
//uf.Error(ufError)
//}
//}

//if ures.Resp == nil {
//uf.Error("No response from server")
//}
//if ures.Resp.StatusCode < 200 ||
//ures.Resp.StatusCode > 299 {
//uf.Error("Bad status code: " + strconv.Itoa(ures.Resp.StatusCode))
//}
//}

//func (ures *UfResponse) AddError(apiError uf.ApiError) {
//if ures.Errors == nil {
//ures.Errors = []uf.ApiError{apiError}
//} else {
//ures.Errors = append(ures.Errors, apiError)
//}
//}

//func (ures *UfResponse) AddErrors(apiErrors []uf.ApiError) {
//if ures.Errors == nil {
//ures.Errors = apiErrors
//} else {
//for _, apiError := range apiErrors {
//ures.Errors = append(ures.Errors, apiError)
//}
//}
//}

//func (ures *UfResponse) Errored() bool {
//if ures.Errors != nil && len(ures.Errors) > 0 {
//return true
//}

//if ures.Resp == nil ||
//ures.Resp.StatusCode < 200 ||
//ures.Resp.StatusCode > 299 {
//return true
//}

//return false
//}

//func (ures *UfResponse) unpackUfResponseErrors(reader io.Reader) {
//respErrors := ApiErrorsFromBody{Errors: []uf.ApiError{}}
//err := json.NewDecoder(reader).Decode(&respErrors)
//if err != nil {
//ures.AddError(uf.DecodeUfResponseErrorsError)
//}

//for _, ApiError := range respErrors.Errors {
//ures.AddError(ApiError)
//}
//}
