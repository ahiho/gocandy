package idgen

import "github.com/bwmarrin/snowflake"

var _node *snowflake.Node

func MustInit(nodeID int64, nodeBits uint8, sequenceBits uint8) error {
	snowflake.NodeBits = nodeBits
	snowflake.StepBits = sequenceBits
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return err
	}
	_node = node
	return nil
}

func GenID() int64 {
	return _node.Generate().Int64()
}

func GenIds(n int) []int64 {
	ids := make([]int64, n)

	for i := 0; i < n; i++ {
		ids[i] = _node.Generate().Int64()
	}
	return ids
}
