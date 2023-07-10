package nodes

import (
	"github.com/gamelee/robot/behavior"
)

func InitNodes() {
	behavior.RegNode(&IF{})
}
