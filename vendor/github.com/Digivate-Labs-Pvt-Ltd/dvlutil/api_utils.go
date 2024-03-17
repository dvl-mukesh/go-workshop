package dvlutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
)

// ApiResponse represents the response from an HTTP request.
type ApiResponse struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int

	// Body contains the raw byte data of the response body.
	Body []byte
}

// BodyToMap converts the body of an ApiResponse into a map[string]interface{}.
// It returns the resulting map and an error if the conversion fails
func (ar *ApiResponse) BodyToMap() (map[string]any, error) {

	body := make(map[string]any)

	err := json.Unmarshal(ar.Body, &body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// DecodeBody decodes the raw byte data in the ApiResponse's Body
// into the specified struct/map (val). It returns an error if the decoding process fails.
func (ar *ApiResponse) DecodeBody(val interface{}) error {
	targetValue := reflect.ValueOf(val)

	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return fmt.Errorf("val must be a non-nil pointer")
	}

	err := json.Unmarshal(ar.Body, targetValue.Interface())

	if err != nil {
		return err
	}

	return nil
}

// CallAPIWithToken sends a POST request to the specified queryURL with the provided token
// in the Authorization header, along with the given request payload.
//
// Parameters:
//   - queryURL: The URL to send the API request to.
//   - token: The authentication token to include in the Authorization header.
//   - request: A map representing the request payload to be sent in the POST request body.
//   - options: Optional configuration options for the HTTP request. These options can be used to configure the HTTP client,
//     such as setting custom headers, timeouts, or other parameters.
//
// Returns:
//   - ApiResponse: A struct containing the API response data.
//   - error: An error, if any, encountered during the API call.
func CallAPIWithToken(queryURL, token string, request map[string]interface{}, options ...HttpOption) (*ApiResponse, error) {

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	return sendPostRequestWithBody(request, queryURL, headers, options...)
}

// CallAPIWithoutToken sends a POST request to the specified queryURL without including any
// authentication token in the headers, along with the given request payload.
//
// Parameters:
//   - queryURL: The URL to send the API request to.
//   - request: A map representing the request payload to be sent in the POST request body.
//   - options: Optional configuration options for the HTTP request. These options can be used to configure the HTTP client,
//     such as setting custom headers, timeouts, or other parameters.
//
// Returns:
//   - ApiResponse: A struct containing the API response data.
//   - error: An error, if any, encountered during the API call.
func CallAPIWithoutToken(queryURL string, request map[string]interface{}, options ...HttpOption) (*ApiResponse, error) {
	return sendPostRequestWithBody(request, queryURL, nil, options...)
}

// GetAuthToken retrieves an authentication token by sending a request to the specified
// authentication URI with the provided user ID and password.
//
// Parameters:
//   - authURI: The URI for obtaining authentication tokens.
//   - userId: The user ID used for authentication.
//   - password: The password associated with the user ID.
//   - options: Optional configuration options for the HTTP request. These options can be used to configure the HTTP client,
//     such as setting custom headers, timeouts, or other parameters.
//
// Returns:
//   - string: The authentication token obtained from the authentication server.
//   - error: An error, if any, encountered during the authentication process.
//
// NOTE:- GetAuthToken is designed to parse the authentication API response specifically for Anota product.
func GetAuthToken(authURI, userId, password string, options ...HttpOption) (string, error) {
	body := new(bytes.Buffer)

	mw := multipart.NewWriter(body)
	mw.WriteField("userid", userId)
	mw.WriteField("password", password)
	mw.Close()

	headers := map[string]string{
		"Content-Type": mw.FormDataContentType(),
	}

	resp, err := sendRequest(authURI, *body, headers, options...)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("auth api returned statusCode: %d, body: %s", resp.StatusCode, resp.Body)
	}

	resM, err := resp.BodyToMap()

	if err != nil {
		return "", err
	}

	value, ok := resM["Value"]

	if !ok {
		return "", fmt.Errorf("unable to fetch token, value not found inside response")
	}

	token, ok := value.(string)

	if !ok {
		return "", fmt.Errorf("value is not a string inside the response")
	}

	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}

// WriteJSON writes a JSON response to the provided http.ResponseWriter with the given status code and value.
//
// The StatusCode parameter sets the HTTP status code for the response.
//
// The val parameter represents the data to be encoded as JSON and sent in the response body. It can be of any type,
// as it will be marshaled to JSON using the encoding/json package.
//
// Note: The caller is responsible for ensuring that no additional content is written to the http.ResponseWriter
// after calling WriteJSON. This function does not terminate the request/response cycle, and it is the caller's
// responsibility to handle any further interactions with the http.ResponseWriter.
func WriteJSON(w http.ResponseWriter, StatusCode int, val any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(StatusCode)
	valJson, _ := json.Marshal(val)
	w.Write(valJson)
}

// GetUnmarshalError returns a descriptive error message when encountering
// a mismatch in field types during the unmarshalling of structs.
//
// The provided 'err' parameter should be an error resulting from
// attempting to unmarshal a JSON or other serialized data into a Go struct.
// This function checks if the error is due to a mismatch in the types of fields
// and returns a human-readable error message with details on the problematic field.
//
// Example:
//
//	type MyStruct struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	data := []byte(`{"name": "John", "age": "30"}`)
//	var result MyStruct
//	err := json.Unmarshal(data, &result)
//	if err != nil {
//	    errorMessage := GetUnmarshalError(err)
//	    fmt.Println(errorMessage)
//	}
//
// In this example, if the 'age' field in the JSON data has a string value instead
// of an integer, GetUnmarshalError will provide a descriptive error message.
//
// If the error is not related to a type mismatch, the function returns
// the error string.
func GetUnmarshalError(err error) string {

	if strings.Contains(err.Error(), "EOF") {
		return INVALID_REQUEST
	}

	msgs := strings.Split(err.Error(), "json: cannot unmarshal ")
	var msg string

	if len(msgs) >= 2 {
		msg = msgs[1]
		msg = "Provided dataType " + strings.Replace(msg, "into Go struct field.", "but required dataType for the key ", -1)
		msg = strings.Replace(msg, "into Go struct field ", "but required dataType for the key ", -1)
	}

	if msg == "" {
		msg = err.Error()
	}
	return msg
}

// Unpanic is a middleware function designed to handle panics that may occur during the execution
// of subsequent HTTP handlers. It takes an http.Handler as input and returns a new http.Handler.
// The returned handler gracefully recovers from panics and responds with a JSON error message
// and a status code of 500 (Internal Server Error) in case of a panic.
func Unpanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("recovered from panic:", r)
				stackTrace := debug.Stack()
				log.Println("StackTrace:", string(stackTrace))
				executablePath, err := os.Executable()
				if err != nil {
					log.Printf("Error getting executable path: %v", err)
				}
				executableName := filepath.Base(executablePath)

				log.Println("Current executing service name:", executableName)

				req := map[string]interface{}{
					"serviceName":  executableName,
					"errorMessage": fmt.Sprintf("%v", r),
					"stackTrace":   string(stackTrace),
				}

				// Print the log request for debugging
				log.Printf("Log request before sending: %+v", req)

				logURL := os.Getenv("SERVICE_PANIC_LOG_URL")
				if logURL != "" {
					res, err := CallAPIWithoutToken(logURL, req)
					if err != nil {
						log.Printf(fmt.Sprintf("%v", err))
					} else if res.StatusCode != http.StatusOK {
						resp, _ := res.BodyToMap()
						log.Printf(fmt.Sprintf("%v", resp))
					}
				} else {
					log.Printf("SERVICE_PANIC_LOG_URL not found in evironment variable")
				}

				WriteJSON(w, http.StatusInternalServerError,
					Response{
						Status: StatusCodeNotOK,
						Msg:    fmt.Sprintf("%v", r),
					})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
