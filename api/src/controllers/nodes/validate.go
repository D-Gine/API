package nodes

import (
	"api/src/database"
	"net/http"

	"github.com/D-Gine/Gestionnaire2Type/validate"
	"github.com/gin-gonic/gin"
)

// @BasePath api/:ruleset_id/node/:id/validate
// Nodes godoc
// @Summary Nodes - Validate one
// @Schemes
// @Description validate the node present the query<br><br>Will return a json object with one key (the node_id) and the error message as value, empty json object if success
// @Tags nodes
// @Accept json
// @Produce json
// @Param creds body validate.NodeList true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router /api/:ruleset_id/node/validate [post]
func ValidateNode(c *gin.Context) {
	var body validate.NodeList
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}

	error_list, err := validate.ValidateNodesTypes(database.Db, body.NodeValues, []string{c.Param("id")}, c.Param("ruleset_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := http.StatusOK
	if len(error_list) != 0 {
		code = http.StatusUnauthorized
	}
	c.JSON(code, error_list)
}

// @BasePath api/:ruleset_id/node/validate
// Nodes godoc
// @Summary Nodes - Validate multiple
// @Schemes
// @Description validate the nodes present in the "node_ids" array in the body<br><br>Will return a json object with every key the invalid nodes and value the error output, empty json object if success
// @Tags nodes
// @Accept json
// @Produce json
// @Param creds body validate.MultipleNodeValidation true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router /api/:ruleset_id/node/:id/validate [post]
func ValidateNodes(c *gin.Context) {
	var body validate.MultipleNodeValidation
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}

	error_list, err := validate.ValidateNodesTypes(database.Db, body.NodeValues, body.NodeIds, c.Param("ruleset_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := http.StatusOK
	if len(error_list) != 0 {
		code = http.StatusUnauthorized
	}
	c.JSON(code, error_list)
}
