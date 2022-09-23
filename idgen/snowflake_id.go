package idgen

import "github.com/bwmarrin/snowflake"

var _node *snowflake.Node

func MustInit(nodeID int64) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		panic(err)
	}
	_node = node
}

func GenID() int64 {
	return _node.Generate().Int64()
}

func GenIds(n int) []int64 {
	ids := make([]int64, 0)
	for i := 0; i < n; i++ {
		ids = append(ids, _node.Generate().Int64())
	}
	return ids
}
