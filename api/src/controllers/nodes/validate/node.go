package validate

import (
	"api/src/database"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func runValidateScript(value json.RawMessage, script string, restrictions json.RawMessage, componentList *ComponentMap) error {
	return nil // IMPLEMENT SCRIPTING CALLING
}

func ValidateNodesTypes(nodeList *NodeList, nodeIds []string, rulesetId string) (map[string]string, error) {
	nodeKeys := slices.Collect(maps.Keys(nodeList.NodeValues))
	nodeRows, err := database.Db.Query(`
		SELECT id, template_id FROM nodes WHERE id = ANY($1)
	`, pq.Array(nodeKeys))
	if err != nil {
		return nil, err
	}
	defer func() { _ = nodeRows.Close() }()
	componentList := make(ComponentMap)
	for nodeRows.Next() {
		var nodeId, compId string
		if err := nodeRows.Scan(&nodeId, &compId); err != nil {
			return nil, err
		}
		if val, exists := nodeList.NodeValues[nodeId]; exists {
			fmt.Println(compId)
			componentList[compId] = val
		}
	}
	if err := nodeRows.Err(); err != nil {
		return nil, err
	}

	verifyRows, err := database.Db.Query(`
		SELECT * FROM get_nodes_verifications($1, $2::uuid)
	`, pq.Array(nodeIds), rulesetId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = verifyRows.Close() }()

	errorMap, count, err := ValidateRowTypes(&componentList, verifyRows)
	if err != nil {
		return errorMap, err
	}

	if count != len(nodeIds) {
		return nil, errors.New("some sent nodes are invalid, cannot compute")
	}

	return errorMap, nil
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

	error_list, err := ValidateNodesTypes(&body, []string{c.Param("id")}, c.Param("ruleset_id"))
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

	error_list, err := ValidateNodesTypes(&body.NodeList, body.NodeIds, c.Param("ruleset_id"))
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
