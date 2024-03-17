package dvlutil

const (
	METHOD_NOT_ALLOWED        = "Method not allowed"
	INVALID_REQUEST           = "Invalid request"
	RECORD_FOUND_SUCCESSFULLY = "Record found successfully"
)

// StatusCodeOK represents the status code indicating a successful operation.
const StatusCodeOK string = "OK"

// StatusCodeNotOK represents the status code indicating a failed or unsuccessful operation.
const StatusCodeNotOK string = "NotOK"

// Response represents a standardized structure for API responses,
// providing a consistent format across different endpoints.
//
// The fields in this struct include:
//   - Status: A string indicating the status of the API response (e.g., "OK" or "NotOK").
//   - Msg: A string containing a human-readable message describing the status or any additional information.
//   - Data: An interface{} field that can hold any type of data as part of the response payload.
//     This allows flexibility in returning various data types such as objects, arrays, or primitive values.
//
// Example usage:
//
//	// Successful response
//	successResponse := Response{
//	    Status: "OK",
//	    Msg:    "Operation completed successfully.",
//	    Data:   someDataStruct,
//	}
//
//	// Error response
//	errorResponse := Response{
//	    Status: "NotOK",
//	    Msg:    "An error occurred during the operation.",
//	    Data:   nil, // Optional: Include additional data related to the error.
//	}
//
//	// Sending the jsonResponse as part of an HTTP response, e.g., using http.ResponseWriter
//	WriteJSON(http.ResponseWriter, http.StatusOK, jsonResponse)
//
// Note: It's important to ensure that the 'Data' field is appropriately typed
// based on the specific API endpoint requirements. The 'Msg' field can be used
// to provide additional context or details about the response.
type Response struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}
