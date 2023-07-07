package behavior

func RegNode(node IBaseNode) {
	defaultRegMap.Register(node)
}

func Nodes() map[string]*Node {
	return defaultRegMap.Nodes()
}
func CreateNode(name string) (interface{}, error) {
	return defaultRegMap.CreateElem(name)
}

func init() {
	RegNode(&Nop{})
	RegNode(&Sequence{})
	RegNode(&Wait{})
	RegNode(&SubTree{})
	RegNode(&RepeaterTimes{})
}
