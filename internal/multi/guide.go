// Package multi - guide tree construction for progressive multiple sequence alignment.
package multi

// GuideNode represents a node in the guide tree (binary tree for MSA ordering).
type GuideNode struct {
	ID       int        // leaf sequence index, or -1 for internal
	Left     *GuideNode
	Right    *GuideNode
	Distance float64 // branch length to parent
}

// IsLeaf reports whether this node is a leaf (represents a single sequence).
func (n *GuideNode) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

// LeafCount returns the number of leaf nodes in the subtree.
func (n *GuideNode) LeafCount() int {
	if n == nil {
		return 0
	}
	if n.IsLeaf() {
		return 1
	}
	return n.Left.LeafCount() + n.Right.LeafCount()
}

// BuildGuideTree constructs a UPGMA guide tree from a distance matrix.
// The matrix should be symmetric with zero diagonal.
func BuildGuideTree(dist [][]float64) *GuideNode {
	n := len(dist)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return &GuideNode{ID: 0}
	}

	// 初始化：每个序列作为一个叶节点
	nodes := make([]*GuideNode, n)
	for i := range nodes {
		nodes[i] = &GuideNode{ID: i}
	}
	sizes := make([]int, n)
	for i := range sizes {
		sizes[i] = 1
	}
	// 复制距离矩阵
	d := make([][]float64, n)
	for i := range d {
		d[i] = make([]float64, n)
		copy(d[i], dist[i])
	}

	active := make([]bool, n)
	for i := range active {
		active[i] = true
	}

	for step := 0; step < n-1; step++ {
		// 找最近的一对
		minD := d[0][1]
		mi, mj := 0, 1
		first := true
		for i := 0; i < len(d); i++ {
			if !active[i] {
				continue
			}
			for j := i + 1; j < len(d); j++ {
				if !active[j] {
					continue
				}
				if first || d[i][j] < minD {
					minD = d[i][j]
					mi, mj = i, j
					first = false
				}
			}
		}
		// 合并 mi 和 mj
		newNode := &GuideNode{
			ID:    -1,
			Left:  nodes[mi],
			Right: nodes[mj],
		}
		nodes[mi].Distance = minD / 2
		nodes[mj].Distance = minD / 2

		// UPGMA 更新距离
		for k := 0; k < len(d); k++ {
			if !active[k] || k == mi || k == mj {
				continue
			}
			newDist := (d[mi][k]*float64(sizes[mi]) + d[mj][k]*float64(sizes[mj])) /
				float64(sizes[mi]+sizes[mj])
			d[mi][k] = newDist
			d[k][mi] = newDist
		}
		sizes[mi] = sizes[mi] + sizes[mj]
		nodes[mi] = newNode
		active[mj] = false
	}

	// 返回最后活跃的节点
	for i, a := range active {
		if a {
			return nodes[i]
		}
	}
	return nil
}
