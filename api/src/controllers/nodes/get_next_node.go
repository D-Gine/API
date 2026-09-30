package nodes

import (
	"api/src/database"
	"api/src/internal/scripting"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var scriptEngine scripting.Engine = scripting.NewStarlarkEngine()

const maxNextNodeDepth = 67000

// @BasePath /api/nodes
// Node godoc
// @Summary Gets next node info for a target node
// @Schemes
// @Description Returns next node metadata for a given target node
// @Tags nodes
// @Accept json
// @Produce json
// @Param ruleset_id path string true "id of the ruleset"
// @Param target_node path string true "id of the target node"
// @Param payload body GetNextNodeRequest true "node values and current tree"
// @Success 200 {object} GetNextNodeResponse
// @Router /api/nodes/{ruleset_id}/get_next_node/{target_node} [post]
func GetNextNode(c *gin.Context) {
	rulesetID := c.Param("ruleset_id")
	targetNodeID := c.Param("target_node")
	if rulesetID == "" || targetNodeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing path params"})
		return
	}

	var body GetNextNodeRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	currentTreeMap := map[string]any{}
	if body.CurrentTree != nil {
		rawCurrentTree, err := json.Marshal(body.CurrentTree)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode current_tree"})
			return
		}
		if err := json.Unmarshal(rawCurrentTree, &currentTreeMap); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse current_tree"})
			return
		}
	}

	targetData, err := loadNodeData(rulesetID, targetNodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Target node not found for this ruleset"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// No script on target = last node
	if !targetData.NextScriptID.Valid || targetData.NextScriptID.String == "" {
		c.JSON(http.StatusOK, GetNextNodeResponse{
			Next: &NextNodeResponse{
				NodeId: targetData.NodeID,
				Final:  true,
				Type:   targetData.Type,
			},
		})
		return
	}

	if !targetData.Script.Valid || targetData.Script.String == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "next script id is set but script content is missing"})
		return
	}

	visited := map[string]struct{}{targetNodeID: {}}
	nextChain, err := resolveNextChain(
		c,
		rulesetID,
		targetData,
		body.NodeValues,
		currentTreeMap,
		visited,
		maxNextNodeDepth,
	)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, errScriptNodeNotFound) || errors.Is(err, errCycleDetected) || errors.Is(err, errDepthExceeded) {
			status = http.StatusBadRequest
		}
		c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetNextNodeResponse{Next: nextChain})
}

var (
	errScriptNodeNotFound = errors.New("script returned unknown next node for this ruleset")
	errCycleDetected      = errors.New("cycle detected while resolving next nodes")
	errDepthExceeded      = errors.New("max next-node depth exceeded")
)

func resolveNextChain(
	c *gin.Context,
	rulesetID string,
	current DbNodeData,
	nodeValues map[string]any,
	currentTree map[string]any,
	visited map[string]struct{},
	remainingDepth int,
) (*NextNodeResponse, error) {
	if remainingDepth <= 0 {
		return nil, errDepthExceeded
	}

	rawResult, err := scriptEngine.Eval(
		c.Request.Context(),
		current.Script.String,
		"get_next_node",
		[]any{nodeValues, currentTree, current.NodeID},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("script execution error: %w", err)
	}

	nextNodeID, err := scripting.AsOptionalString(rawResult)
	if err != nil {
		return nil, fmt.Errorf("script return type error: %w", err)
	}
	if nextNodeID == nil {
		return nil, nil
	}

	if _, alreadySeen := visited[*nextNodeID]; alreadySeen {
		return nil, errCycleDetected
	}
	visited[*nextNodeID] = struct{}{}

	nextData, err := loadNodeData(rulesetID, *nextNodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errScriptNodeNotFound
		}
		return nil, err
	}

	response := &NextNodeResponse{
		NodeId: nextData.NodeID,
		Type:   nextData.Type,
		Final:  !nextData.NextScriptID.Valid || nextData.NextScriptID.String == "",
	}

	if response.Final {
		return response, nil
	}

	if !nextData.Script.Valid || nextData.Script.String == "" {
		return nil, fmt.Errorf("next script id is set but script content is missing for node %s", nextData.NodeID)
	}

	appendNodeToTree(currentTree, nextData.NodeID)
	child, err := resolveNextChain(c, rulesetID, nextData, nodeValues, currentTree, visited, remainingDepth-1)
	if err != nil {
		return nil, err
	}
	response.Next = child
	return response, nil
}

func loadNodeData(rulesetID string, nodeID string) (DbNodeData, error) {
	var out DbNodeData
	var typeRaw []byte

	err := database.Db.QueryRow(`
		SELECT
			n.id::text,
			n.next_script_id::text,
			ns.script,
			COALESCE(get_base_type_json(sct.type_id), '{}'::jsonb)::text
		FROM public.nodes n
		INNER JOIN public.static_components_templates sct
			ON sct.id = n.template_id
		LEFT JOIN public.next_scripts ns
			ON ns.id = n.next_script_id
		WHERE n.ruleset_id = $1 AND n.id = $2
		LIMIT 1
	`, rulesetID, nodeID).Scan(&out.NodeID, &out.NextScriptID, &out.Script, &typeRaw)
	if err != nil {
		return out, err
	}

	if len(typeRaw) == 0 {
		typeRaw = []byte(`{}`)
	}
	if err := json.Unmarshal(typeRaw, &out.Type); err != nil {
		return out, fmt.Errorf("invalid type payload from database for node %s", nodeID)
	}

	return out, nil
}

func appendNodeToTree(tree map[string]any, nodeID string) {
	if tree == nil {
		return
	}
	if len(tree) == 0 {
		tree["node_id"] = nodeID
		return
	}

	current := tree
	for {
		if id, ok := current["node_id"].(string); ok && id == nodeID {
			return
		}

		nextAny, hasNext := current["next"]
		if !hasNext || nextAny == nil {
			current["next"] = map[string]any{"node_id": nodeID}
			return
		}

		nextMap, ok := nextAny.(map[string]any)
		if !ok {
			current["next"] = map[string]any{"node_id": nodeID}
			return
		}
		current = nextMap
	}
}
