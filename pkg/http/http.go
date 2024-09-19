package http

import (
	"encoding/json"
	"fmt"

	"github.com/KeeganBeuthin/TBV-Go-SDK/pkg/utils"
)

type Request struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func ParseRequest(requestPtr *byte) Request {
	requestStr := utils.PtrToString(requestPtr)
	var request Request
	err := json.Unmarshal([]byte(requestStr), &request)
	if err != nil {
		fmt.Printf("Error parsing request: %v\n", err)
	}
	return request
}

func (r Response) Stringify() string {
	jsonResponse, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("Error stringifying response: %v\n", err)
		return ""
	}
	return string(jsonResponse)
}

func handle_http_request(requestPtr *byte) *byte {
	request := ParseRequest(requestPtr)
	response := HandleRequest(request)
	return utils.StringToPtr(response.Stringify())
}

func HandleRequest(req Request) Response {
	switch req.Path {
	case "/api/data":
		return handleDataRequest(req)
	default:
		return Response{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "Not Found",
		}
	}
}

func handleDataRequest(req Request) Response {
	switch req.Method {
	case "GET":
		data := map[string]string{"message": "Hello from WebAssembly API!"}
		jsonData, _ := json.Marshal(data)
		return Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(jsonData),
		}
	case "POST":
		return Response{
			StatusCode: 201,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"message": "Data created successfully"}`,
		}
	case "PUT":
		return Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"message": "Data updated successfully"}`,
		}
	case "DELETE":
		return Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"message": "Data deleted successfully"}`,
		}
	default:
		return Response{
			StatusCode: 405,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "Method Not Allowed",
		}
	}
}
