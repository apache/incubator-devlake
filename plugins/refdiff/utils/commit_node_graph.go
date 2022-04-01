package utils

type CommitNode struct {
	Sha    string
	Parent []*CommitNode
}

type CommitNodeGraph struct {
	node map[string]*CommitNode
}

func NewCommitNodeGraph() *CommitNodeGraph {
	return &CommitNodeGraph{
		node: make(map[string]*CommitNode),
	}
}

func (cng *CommitNodeGraph) AddSide(commit_sha string, commit_parent_sha string) {
	var CommitNow *CommitNode
	var CommitParentNow *CommitNode
	var ok bool
	if CommitNow, ok = cng.node[commit_sha]; !ok {
		CommitNow = &CommitNode{
			Sha: commit_sha,
		}
		cng.node[commit_sha] = CommitNow
	}

	if CommitParentNow, ok = cng.node[commit_parent_sha]; !ok {
		CommitParentNow = &CommitNode{
			Sha: commit_parent_sha,
		}
		cng.node[commit_parent_sha] = CommitParentNow
	}
	CommitNow.Parent = append(CommitNow.Parent, CommitParentNow)
}

func (cng *CommitNodeGraph) CalculateLost(source_sha string, target_sha string) ([]string, int, int) {
	var oldCommitNode *CommitNode
	var newCommitNode *CommitNode
	var lostSha []string
	var ok bool
	if oldCommitNode, ok = cng.node[source_sha]; !ok {
		oldCommitNode = &CommitNode{
			Sha: source_sha,
		}
	}
	if newCommitNode, ok = cng.node[target_sha]; !ok {
		newCommitNode = &CommitNode{
			Sha: target_sha,
		}
	}

	var oldGroup map[string]*CommitNode = make(map[string]*CommitNode)
	type funcType func(*CommitNode)
	var dfs funcType
	dfs = func(now *CommitNode) {
		if _, ok = oldGroup[now.Sha]; ok {
			return
		}
		oldGroup[now.Sha] = now
		for _, node := range now.Parent {
			dfs(node)
		}
	}
	dfs(oldCommitNode)

	var newGroup map[string]*CommitNode = make(map[string]*CommitNode)
	dfs = func(now *CommitNode) {
		if _, ok = oldGroup[now.Sha]; ok {
			return
		}
		if _, ok = newGroup[now.Sha]; ok {
			return
		}
		newGroup[now.Sha] = now
		lostSha = append(lostSha, now.Sha)
		for _, node := range now.Parent {
			dfs(node)
		}
	}
	dfs(newCommitNode)

	return lostSha, len(oldGroup), len(newGroup)
}

func (cng *CommitNodeGraph) Size() int {
	return len(cng.node)
}
