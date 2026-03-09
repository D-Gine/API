/*
** D&GINE Project, 2026
** API
** File description:
** characters/creation/next_node.go
 */

package charactersCreate

import (
	"api/src/database"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type NextNodeArgs struct {
	NodeId string `json:"node_id"` // actual node id
	Type   string `json:"type"`
	Value  string `json:"value"`
}

type row struct {
	ChildId    string `json:"child_id"`
	Name       string `json:"name"`
	Expression string `json:"expression"`
	Type       string `json:"type"`
	Metadata   string `json:"metadata"`
}

type NextNodeResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Metadata string `json:"metadata"`
}

func valueMatchCondition(value string, valType string, expression string, expType string) bool {
	if valType != expType {
		return false
	}

	var operator, operand string
	for _, op := range []string{"<=", ">=", "!=", "<", ">", "=="} {
		if strings.HasPrefix(expression, op) {
			operator = op
			operand = strings.TrimSpace(expression[len(op):])
			break
		}
	}
	if operator == "" {
		return false // no operator ? or return true maybe @Jands
	}

	switch valType { // modify with enums of types set in database script @Jands
	case "int":
		v, err1 := strconv.Atoi(value)
		e, err2 := strconv.Atoi(operand)
		if err1 != nil || err2 != nil {
			return false
		}
		switch operator {
		case "<":
			return v < e
		case "<=":
			return v <= e
		case ">":
			return v > e
		case ">=":
			return v >= e
		case "==":
			return v == e
		case "!=":
			return v != e
		}

	case "float":
		v, err1 := strconv.ParseFloat(value, 64)
		e, err2 := strconv.ParseFloat(operand, 64)
		if err1 != nil || err2 != nil {
			return false
		}
		switch operator {
		case "<":
			return v < e
		case "<=":
			return v <= e
		case ">":
			return v > e
		case ">=":
			return v >= e
		case "==":
			return v == e
		case "!=":
			return v != e
		}

	case "string":
		switch operator {
		case "==":
			return value == operand
		case "!=":
			return value != operand
		default:
			return false // ordering not supported for strings
		}
	}

	return false
}

// @BasePath /api/characters/create/nextnode
// Characters godoc
// @Summary Returns the next node
// @Schemes
// @Description Returns the next node of the given node depending on the ruleset character creation tree<br><br>This will check for the tree for all childs of the given node, and return the one(s) that match the condition with the given value<br><br>If no child matches the condition, it means that the character creation process is finished, and the frontend can send the character creation submission request
// @Tags characters creation
// @Accept json
// @Produce json
// @Param creds body NextNodeArgs true "Actual node id, along with the value and type of the component"
// @Success 200 {array} NextNodeResponse
// @Router /api/characters/create/nextnode [get]
func CharacterNextNode(c *gin.Context) {
	var args NextNodeArgs

	result := []NextNodeResponse{}
	// Body Json parsing
	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	childs := []row{}
	rows, err := database.Db.Query(
		`SELECT
			cn.child_id, ct.name, ct.expression, ctype.type, ctype.metadata
		FROM games.creation_nodes cn
		INNER JOIN games.components_templates ct
			ON cn.child_id=ct.id
		INNER JOIN games.component_type ctype
				ON ct.type=ctype.id
		WHERE cn.parent_id=$1;`, args.NodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusOK, result) // the node was the last
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error : " + err.Error()})
		return
	}

	defer rows.Close()
	for rows.Next() {
		var childId, name, expression, valType, metadata sql.NullString
		rows.Scan(&childId, &name, &expression, &valType, &metadata)
		childs = append(childs, row{
			ChildId:    childId.String,
			Name:       name.String,
			Expression: expression.String,
			Type:       valType.String,
			Metadata:   metadata.String,
		})
	}

	for row := range childs {
		if childs[row].Expression == "" { // no expression, we add all childs @Jands
			result = append(result, NextNodeResponse{
				Id:       childs[row].ChildId,
				Name:     childs[row].Name,
				Type:     childs[row].Type,
				Metadata: childs[row].Metadata,
			})
		} else if valueMatchCondition(args.Value, args.Type, childs[row].Expression, childs[row].Type) {
			result = append(result, NextNodeResponse{
				Id:       childs[row].ChildId,
				Name:     childs[row].Name,
				Type:     childs[row].Type,
				Metadata: childs[row].Metadata,
			})
		}
	}

	c.JSON(http.StatusOK, result)
}
