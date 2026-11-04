package apidata

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/gin-gonic/gin"
)

var scalarConfig = map[string]any{
	"theme":                           "saturn",
	"layout":                          "modern",
	"orderSchemaPropertiesBy":         "preserve",
	"orderRequiredPropertiesFirst":    true,
	"showOperationId":                 true,
	"operationTitleSource":            "summary",
	"hideTestRequestButton":           false,
	"hideClientButton":                true,
	"showSidebar":                     true,
	"showDeveloperTools":              "never",
	"showToolbar":                     "localhost",
	"persistAuth":                     false,
	"telemetry":                       false,
	"hideModels":                      false,
	"hideModelNames":                  false,
	"documentDownloadType":            "both",
	"hideSearch":                      false,
	"hideDarkModeToggle":              false,
	"withDefaultFonts":                true,
	"defaultOpenFirstTag":             true,
	"defaultOpenAllTags":              false,
	"expandAllModelSections":          false,
	"expandAllParameters":             true,
	"expandAllResponses":              false,
	"expandAllSchemaProperties":       false,
	"maxVisibleRequestBodyProperties": 12,
	"schemaKeyboardNav":               false,
	"modelsSectionLabel":              "Models",
	"agent": map[string]any{
		"disabled": true,
	},
	"featuredClients": []string{
		"javascript/fetch",
		"shell/curl",
		"node/fetch",
	},
}

func ScalarDoc(c *gin.Context) {
	cfg, err := json.Marshal(scalarConfig)
	if err != nil {
		c.String(http.StatusInternalServerError, "invalid scalar config")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <title>D&GINE API</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/swagger/doc.json"
      data-configuration="%s"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`, html.EscapeString(string(cfg))))
}
