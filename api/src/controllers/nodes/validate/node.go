package validate

import (
	"api/src/database"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func runValidateScript(nodeId string, script string, restrictions json.RawMessage, nodelist *NodeList) error {
	return nil // IMPLEMENT SCRIPTING CALLING
}

func ValidateNodesTypes(nodeList *NodeList, nodeIds []string, rulesetId string) (error, map[string]string) {
	rows, err := database.Db.Query(`
		SELECT * FROM get_nodes_verifications($1, $2::uuid)`, pq.Array(nodeIds), rulesetId,
	)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	var count int
	var errors_map = map[string]string{}
	var result Verifications
	for count = 0; rows.Next(); count++ {
		if err := rows.Scan(&result.NodeId, &result.BaseName, pq.Array(&result.Scripts), &result.Restrictions); err != nil {
			return err, nil
		}
		if err := ValidateBaseType(nodeList.NodeValues[result.NodeId], result.BaseName, result.Restrictions); err != nil {
			errors_map[result.NodeId] = "error validation base type: " + err.Error()
			continue
		}

		for i := 0; i < len(result.Scripts); i++ {
			if err = runValidateScript(
				result.NodeId,
				result.Scripts[i],
				result.Restrictions,
				nodeList,
			); err != nil {
				errors_map[result.NodeId] = err.Error()
				continue
			}
		}

	}
	if err = rows.Err(); err != nil {
		return err, nil
	}

	if count != len(nodeIds) {
		return errors.New("some sent nodes are invalids, cannot compute"), nil
	}
	return nil, errors_map
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

	err, error_list := ValidateNodesTypes(&body, []string{c.Param("id")}, c.Param("ruleset_id"))
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

	err, error_list := ValidateNodesTypes(&body.NodeList, body.NodeIds, c.Param("ruleset_id"))
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
