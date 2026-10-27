package validate

import (
	"api/src/controllers"
	"api/src/database"
	"api/src/internal/scripting"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func runValidateScript(ctx *gin.Context, value json.RawMessage, script string, restrictions json.RawMessage, componentList *controllers.ComponentMap) error {
	rawResult, err := scripting.ScriptEngine.Eval(
		ctx.Request.Context(),
		script,
		"validate_type",
		[]any{value, restrictions, *componentList},
		nil,
	)
	if err != nil {
		return fmt.Errorf("script execution error: %w", err)
	}
	stringError, err := scripting.AsOptionalString(rawResult)
	if err != nil {
		return fmt.Errorf("script return type error: %w", err)
	}
	if stringError == nil {
		return nil
	}
	return errors.New(*stringError)
}

func ValidateNodesTypes(ctx *gin.Context, nodeList *NodeList, nodeIds []string, rulesetId string) (error, map[string]string) {
	nodeKeys := slices.Collect(maps.Keys(nodeList.NodeValues))
	nodeRows, err := database.Db.Query(`
		SELECT id, template_id FROM nodes WHERE id = ANY($1)
	`, pq.Array(nodeKeys))
	if err != nil {
		return err, nil
	}
	defer nodeRows.Close()
	componentList := make(controllers.ComponentMap)
	nodeRecoverMap := make(map[string]string)
	for nodeRows.Next() {
		var nodeId, compId string
		if err := nodeRows.Scan(&nodeId, &compId); err != nil {
			return err, nil
		}
		if val, exists := nodeList.NodeValues[nodeId]; exists {
			componentList[compId] = val
			nodeRecoverMap[compId] = nodeId
		}
	}
	if err := nodeRows.Err(); err != nil {
		return err, nil
	}

	verifyRows, err := database.Db.Query(`
		SELECT * FROM get_nodes_verifications($1, $2::uuid)
	`, pq.Array(nodeIds), rulesetId)
	if err != nil {
		return err, nil
	}
	defer verifyRows.Close()

	err, errorMap, count := ValidateRowTypes(ctx, &componentList, verifyRows)
	nodeErrorMap := make(map[string]string)
	for id, v := range errorMap {
		nodeErrorMap[nodeRecoverMap[id]] = v
	}
	if err != nil {
		return err, nodeErrorMap
	}

	if count != len(nodeIds) {
		return errors.New("some sent nodes are invalid, cannot compute"), nil
	}

	return nil, nodeErrorMap
}

// @BasePath api/:ruleset_id/node/:id/validate
// Nodes godoc
// @Summary validate one node
// @Schemes
// @Description validate the node present the query<br><br>Will return a json object with one key (the node_id) and the error message as value, empty json object if success
// @Tags nodes
// @Accept json
// @Produce json
// @Param creds body CreateRulesetsArgs true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router api/:ruleset_id/node/validate [post]
func ValidateNode(c *gin.Context) {
	var body NodeList
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}

	err, error_list := ValidateNodesTypes(c, &body, []string{c.Param("id")}, c.Param("ruleset_id"))
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
// @Summary validate multiple nodes
// @Schemes
// @Description validate the nodes present in the "node_ids" array in the body<br><br>Will return a json object with every key the invalid nodes and value the error output, empty json object if success
// @Tags nodes
// @Accept json
// @Produce json
// @Param creds body CreateRulesetsArgs true "ruleset related informations that will be later needed for the login process"
// @Success 200 {object} structs.PostResponse
// @Router api/:ruleset_id/node/:id/validate [post]
func ValidateNodes(c *gin.Context) {
	var body MultipleNodeValidation
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Invalid body"})
		return
	}

	err, error_list := ValidateNodesTypes(c, &body.NodeList, body.NodeIds, c.Param("ruleset_id"))
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
