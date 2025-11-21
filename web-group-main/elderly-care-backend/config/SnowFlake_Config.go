package config

import (
	"elderly-care-backend/global"
	"fmt"
	"github.com/bwmarrin/snowflake"
)

func initSnowFlake() {
	node, err := snowflake.NewNode(Config.Server.NodeId)
	if err != nil {
		fmt.Fprintf(logFile, "Cannot init snowflake: %v\n", err)
	}
	global.SnowFlake = node
}
