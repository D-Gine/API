/*
** D&GINE Project, 2026
** API
** File description:
** CRUDTEMPLATE/read_id.go
 */

package dev

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type CallEngineArgs struct {
	Method  string            `json:"method"`
	Uri     string            `json:"uri"`
	Headers map[string]string `json:"headers"`
	Body    map[string]string `json:"body"`
}

type Response struct {
	Code int            `json:"code"`
	Data map[string]any `json:"data"`
}

func httpRequest(method string, uri string, headers map[string]string, body map[string]string, debug bool) (Response, error) {

	// if debug {
	// 	fmt.Printf(`
	// 		=========== HTTP Request ==========
	// 		Method : %s
	// 		Uri : %s
	// 		Header : %v
	// 		Body : %v
	// 		=========== HTTP Request ==========
	// 	`, method, uri, headers, body)
	// }
	data := url.Values{}
	for key, value := range body {
		data.Set(key, value)
	}

	req, err := http.NewRequest(method, uri, strings.NewReader(data.Encode()))
	if err != nil {
		return Response{}, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var result map[string]any
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return Response{}, err
	}

	// if debug {
	// 	fmt.Printf(`
	// 		=========== HTTP Response ==========
	// 		Code : %d
	// 		Body : %v
	// 		=========== HTTP Response ==========
	// 	`, resp.StatusCode, result)
	// }

	return Response{
		Code: resp.StatusCode,
		Data: result,
	}, nil
}

// @BasePath /api/dev/callengine
// CRUDTEMPLATE godoc
// @Summary Reads given CRUDTEMPLATE data
// @Schemes
// @Description ⚠️ Only accessible to admins ⚠️<br><br>Reads the given CRUDTEMPLATE's account data
// @Tags CRUDTEMPLATEs
// @Param id query string true "id of the CRUDTEMPLATE to read"
// @Produce json
// @Success 200 {object} DbReadResponse
// @Router /api/CRUDTEMPLATEs/id [get]
func CallEngine(c *gin.Context) {
	var args CallEngineArgs

	// Body Json parsing
	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := httpRequest(args.Method, args.Uri, args.Headers, args.Body, (os.Getenv("GIN_MODE") == "debug"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
