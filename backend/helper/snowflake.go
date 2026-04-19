package helper

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

var snowflakeNode *snowflake.Node

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(fmt.Sprintf("snowflake node init failed: %v", err))
	}
	snowflakeNode = node
}

func GenerateSnowflake() snowflake.ID {
	return snowflakeNode.Generate()
}
