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

func HandleRequest(req Request) Response {
	switch req.Path {
	case "/html":
		return handleHtmlRequest(req)
	default:
		return Response{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "Not Found",
		}
	}
}

func handleHtmlRequest(req Request) Response {
	switch req.Method {
	case "GET":
		return Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "text/html"},
			Body:       "<html><body><h1>Hello from WebAssembly!</h1></body></html>",
		}
	case "PUT":
		// Update HTML content
		// In a real implementation, you'd store this somewhere
		return Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "HTML content updated",
		}
	case "DELETE":
		// Delete HTML content
		return Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "HTML content deleted",
		}
	default:
		return Response{
			StatusCode: 405,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "Method Not Allowed",
		}
	}
}
