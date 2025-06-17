package stats_collector

import (
	"github.com/thlaurentino/arit/internal/reader"
)

type FunctionStats struct {
	FunctionName                  string
	LinesOfCode                   int
	ParameterCount                int
	MaxNestingDepth               int
	MaxMessageChain               int
	MaxConsecutivePrimitiveParams int
}

func Collect(rootNode *reader.RichNode) []FunctionStats {
	var stats []FunctionStats
	var findFunctions func(node *reader.RichNode)

	findFunctions = func(node *reader.RichNode) {
		if node == nil {
			return
		}

		if isFunctionDefinition(node) {
			stats = append(stats, analyzeFunctionNode(node))
		}

		if !isFunctionDefinition(node) {
			for _, child := range node.Children {
				findFunctions(child)
			}
		}
	}

	findFunctions(rootNode)
	return stats
}

func isFunctionDefinition(node *reader.RichNode) bool {
	if node.Type != reader.NodeList || len(node.Children) < 2 {
		return false
	}

	return node.Children[0].Type == reader.NodeSymbol && (node.Children[0].Value == "defn" || node.Children[0].Value == "defn-")
}

func analyzeFunctionNode(fnNode *reader.RichNode) FunctionStats {
	return FunctionStats{
		FunctionName:                  getFunctionName(fnNode),
		LinesOfCode:                   calculateLinesOfCode(fnNode),
		ParameterCount:                countParameters(fnNode),
		MaxNestingDepth:               calculateNestingForBody(fnNode),
		MaxMessageChain:               calculateMaxMessageChain(fnNode),
		MaxConsecutivePrimitiveParams: countMaxConsecutivePrimitives(fnNode),
	}
}

func calculateLinesOfCode(fnNode *reader.RichNode) int {
	if fnNode.Location == nil || len(fnNode.Children) == 0 {
		return 0
	}

	lastChild := fnNode.Children[len(fnNode.Children)-1]
	if lastChild.Location == nil {

		return fnNode.Location.EndLine - fnNode.Location.StartLine + 1
	}

	return lastChild.Location.EndLine - fnNode.Location.StartLine + 1
}

func countParameters(fnNode *reader.RichNode) int {

	paramCandidateIndex := 2
	if len(fnNode.Children) > paramCandidateIndex && fnNode.Children[paramCandidateIndex].Type == reader.NodeString {
		paramCandidateIndex++
	}
	if len(fnNode.Children) > paramCandidateIndex && fnNode.Children[paramCandidateIndex].Type == reader.NodeMap {
		paramCandidateIndex++
	}

	if len(fnNode.Children) <= paramCandidateIndex {
		return 0
	}

	paramsNode := fnNode.Children[paramCandidateIndex]

	if paramsNode.Type == reader.NodeVector {
		return len(paramsNode.Children)
	}

	if paramsNode.Type == reader.NodeList {
		maxParams := 0
		for _, arityForm := range paramsNode.Children {

			if arityForm.Type == reader.NodeList && len(arityForm.Children) > 0 && arityForm.Children[0].Type == reader.NodeVector {
				numParams := len(arityForm.Children[0].Children)
				if numParams > maxParams {
					maxParams = numParams
				}
			}
		}
		return maxParams
	}

	return 0
}

func calculateNestingForBody(node *reader.RichNode) int {
	if node == nil {
		return 0
	}

	maxOverallCallStackDepth := 0

	var findDepth func(*reader.RichNode) int
	findDepth = func(n *reader.RichNode) int {
		if n == nil {
			return 0
		}

		maxChildCallDepth := 0
		isCall := false

		if n.Type == reader.NodeList && len(n.Children) > 0 {
			firstChildType := n.Children[0].Type
			if firstChildType == reader.NodeSymbol || firstChildType == reader.NodeKeyword {
				isCall = true
			}
		} else if n.Type == reader.NodeFnLiteral {
			isCall = false
			for _, child := range n.Children {
				depth := findDepth(child)
				if depth > maxChildCallDepth {
					maxChildCallDepth = depth
				}
			}
			return maxChildCallDepth
		}

		if isCall {
			for i := 1; i < len(n.Children); i++ {
				argDepth := findDepth(n.Children[i])
				if argDepth > maxChildCallDepth {
					maxChildCallDepth = argDepth
				}
			}
			return 1 + maxChildCallDepth
		} else if n.Type == reader.NodeList || n.Type == reader.NodeVector || n.Type == reader.NodeMap {
			for _, child := range n.Children {
				childDepth := findDepth(child)
				if childDepth > maxChildCallDepth {
					maxChildCallDepth = childDepth
				}
			}
			return maxChildCallDepth
		}

		return 0
	}

	for _, child := range node.Children {
		depth := findDepth(child)
		if depth > maxOverallCallStackDepth {
			maxOverallCallStackDepth = depth
		}
	}

	return maxOverallCallStackDepth
}

func calculateMaxMessageChain(fnNode *reader.RichNode) int {
	maxChain := 0

	var findChains func(node *reader.RichNode)
	findChains = func(node *reader.RichNode) {
		if node == nil {
			return
		}

		if node.Type == reader.NodeList && len(node.Children) >= 3 &&
			node.Children[0].Type == reader.NodeSymbol && node.Children[0].Value == "get-in" &&
			node.Children[2].Type == reader.NodeVector {
			pathVector := node.Children[2]
			if len(pathVector.Children) > maxChain {
				maxChain = len(pathVector.Children)
			}
		}

		if node.Type == reader.NodeList && len(node.Children) > 1 &&
			node.Children[0].Type == reader.NodeSymbol && (node.Children[0].Value == "->" || node.Children[0].Value == "some->") {

			chainLen := len(node.Children) - 2
			if chainLen > maxChain {
				maxChain = chainLen
			}
		}

		for _, child := range node.Children {
			findChains(child)
		}
	}

	findChains(fnNode)
	return maxChain
}

func countMaxConsecutivePrimitives(fnNode *reader.RichNode) int {
	paramsNode := getParamsNode(fnNode)
	if paramsNode == nil {
		return 0
	}

	maxConsecutive := 0
	currentConsecutive := 0

	for _, param := range paramsNode.Children {
		if isPrimitiveLike(param) {
			currentConsecutive++
		} else {
			if currentConsecutive > maxConsecutive {
				maxConsecutive = currentConsecutive
			}
			currentConsecutive = 0
		}
	}
	if currentConsecutive > maxConsecutive {
		maxConsecutive = currentConsecutive
	}

	return maxConsecutive
}

func getParamsNode(fnNode *reader.RichNode) *reader.RichNode {

	paramCandidateIndex := 2
	if len(fnNode.Children) > paramCandidateIndex && fnNode.Children[paramCandidateIndex].Type == reader.NodeString {
		paramCandidateIndex++
	}
	if len(fnNode.Children) > paramCandidateIndex && fnNode.Children[paramCandidateIndex].Type == reader.NodeMap {
		paramCandidateIndex++
	}

	if len(fnNode.Children) > paramCandidateIndex {
		paramsNode := fnNode.Children[paramCandidateIndex]
		if paramsNode.Type == reader.NodeVector {
			return paramsNode
		}
		if paramsNode.Type == reader.NodeList {
			for _, arityForm := range paramsNode.Children {
				if arityForm.Type == reader.NodeList && len(arityForm.Children) > 0 && arityForm.Children[0].Type == reader.NodeVector {
					return arityForm.Children[0]
				}
			}
		}
	}
	return nil
}

func isPrimitiveLike(paramNode *reader.RichNode) bool {
	if paramNode.Type == reader.NodeSymbol {

		val := paramNode.Value
		if val == "&" || val == "_" {
			return false
		}
		return true
	}
	return false
}

func getFunctionName(fnNode *reader.RichNode) string {
	if len(fnNode.Children) < 1 {
		return "<unknown>"
	}
	firstChild := fnNode.Children[0]
	if firstChild.Type != reader.NodeSymbol {
		return "<unknown>"
	}

	if firstChild.Value == "fn" {
		return "<anonymous>"
	}

	if (firstChild.Value == "defn" || firstChild.Value == "defn-") && len(fnNode.Children) > 1 {
		secondChild := fnNode.Children[1]
		if secondChild.Type == reader.NodeSymbol {
			return secondChild.Value
		}
	}

	return "<unknown>"
}
