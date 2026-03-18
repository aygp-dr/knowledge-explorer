package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Node represents a node in the knowledge tree.
type Node struct {
	Name        string
	Description string
	Resources   []string
	Children    []*Node
	Parent      *Node
	Expanded    bool
}

// IsLeaf returns true if the node has no children.
func (n *Node) IsLeaf() bool {
	return len(n.Children) == 0
}

// Depth returns the depth of this node in the tree.
func (n *Node) Depth() int {
	d := 0
	p := n.Parent
	for p != nil {
		d++
		p = p.Parent
	}
	return d
}

// Breadcrumb returns the path from root to this node.
func (n *Node) Breadcrumb() string {
	parts := []string{}
	cur := n
	for cur != nil {
		parts = append([]string{cur.Name}, parts...)
		cur = cur.Parent
	}
	return strings.Join(parts, " > ")
}

// FlattenVisible returns the visible nodes in DFS order (respecting expand/collapse).
func FlattenVisible(roots []*Node) []*Node {
	var result []*Node
	var walk func(nodes []*Node)
	walk = func(nodes []*Node) {
		for _, n := range nodes {
			result = append(result, n)
			if n.Expanded && len(n.Children) > 0 {
				walk(n.Children)
			}
		}
	}
	walk(roots)
	return result
}

// setParents recursively sets parent pointers.
func setParents(nodes []*Node, parent *Node) {
	for _, n := range nodes {
		n.Parent = parent
		setParents(n.Children, n)
	}
}

// BuildKnowledgeTree constructs the CS knowledge tree.
func BuildKnowledgeTree() []*Node {
	tree := []*Node{
		{
			Name:        "Algorithms",
			Description: "Step-by-step procedures for solving computational problems. Algorithms are the foundation of computer science and software engineering.",
			Resources:   []string{"Introduction to Algorithms (CLRS)", "Algorithm Design Manual (Skiena)", "https://visualgo.net"},
			Children: []*Node{
				{
					Name:        "Sorting",
					Description: "Algorithms that arrange elements in a specific order. Key concepts include comparison-based sorts (O(n log n) lower bound), stable vs unstable sorts, and in-place vs out-of-place.",
					Resources:   []string{"Quicksort (Hoare, 1961)", "Mergesort, Heapsort, Timsort", "https://sorting.at"},
				},
				{
					Name:        "Searching",
					Description: "Algorithms for finding elements in data structures. Includes linear search, binary search, interpolation search, and search in graphs/trees.",
					Resources:   []string{"Binary Search variations", "Ternary Search", "https://cp-algorithms.com"},
				},
				{
					Name:        "Graph",
					Description: "Algorithms operating on graph structures: BFS, DFS, shortest path (Dijkstra, Bellman-Ford), minimum spanning tree (Kruskal, Prim), topological sort, and network flow.",
					Resources:   []string{"Graph Theory (Diestel)", "Network Flows (Ahuja et al.)", "https://graphonline.ru/en/"},
				},
				{
					Name:        "Dynamic Programming",
					Description: "Technique for solving problems by breaking them into overlapping subproblems and storing solutions. Key patterns: memoization, tabulation, state compression.",
					Resources:   []string{"Dynamic Programming for Coding Interviews", "Competitive Programmer's Handbook", "https://atcoder.jp/contests/dp"},
				},
			},
		},
		{
			Name:        "Data Structures",
			Description: "Organized formats for storing and accessing data efficiently. The choice of data structure directly impacts algorithm performance and system design.",
			Resources:   []string{"Data Structures and Algorithms in Java (Goodrich)", "Open Data Structures", "https://www.cs.usfca.edu/~galles/visualization/"},
			Children: []*Node{
				{
					Name:        "Arrays",
					Description: "Contiguous memory blocks providing O(1) random access. Variants include dynamic arrays, circular buffers, and bit arrays. Foundation for many other structures.",
					Resources:   []string{"Array-based data structures", "Cache-friendly programming", "https://en.cppreference.com/w/cpp/container/vector"},
				},
				{
					Name:        "Trees",
					Description: "Hierarchical structures with parent-child relationships. Includes binary trees, BSTs, AVL trees, red-black trees, B-trees, tries, and segment trees.",
					Resources:   []string{"Introduction to Algorithms (CLRS) Ch. 12-13", "B-Trees and Database Indexing", "https://visualgo.net/en/bst"},
				},
				{
					Name:        "Graphs",
					Description: "Nodes connected by edges representing relationships. Representations: adjacency matrix, adjacency list, edge list. Used in social networks, maps, and dependency resolution.",
					Resources:   []string{"Graph Theory (Bondy & Murty)", "NetworkX documentation", "https://d3js.org (graph visualization)"},
				},
				{
					Name:        "HashMaps",
					Description: "Key-value stores with average O(1) lookup using hash functions. Collision resolution: chaining, open addressing. Used in caches, symbol tables, and deduplication.",
					Resources:   []string{"Hash Table internals", "Consistent Hashing", "https://opendatastructures.org/ods-java/5_Hash_Tables.html"},
				},
			},
		},
		{
			Name:        "Systems",
			Description: "Design and implementation of computing systems at various scales. Covers operating systems, networking, distributed systems, and data management.",
			Resources:   []string{"Designing Data-Intensive Applications (Kleppmann)", "Computer Systems: A Programmer's Perspective", "https://teachyourselfcs.com"},
			Children: []*Node{
				{
					Name:        "OS",
					Description: "Operating systems manage hardware resources and provide abstractions: processes, threads, memory management, file systems, scheduling, and system calls.",
					Resources:   []string{"Operating Systems: Three Easy Pieces (OSTEP)", "Linux Kernel Development (Love)", "https://pages.cs.wisc.edu/~remzi/OSTEP/"},
				},
				{
					Name:        "Networking",
					Description: "Communication between computing devices. Covers the OSI model, TCP/IP stack, HTTP/HTTPS, DNS, load balancing, CDNs, and network security.",
					Resources:   []string{"Computer Networking: A Top-Down Approach (Kurose)", "TCP/IP Illustrated (Stevens)", "https://beej.us/guide/bgnet/"},
				},
				{
					Name:        "Distributed",
					Description: "Systems spanning multiple machines: consensus (Raft, Paxos), replication, partitioning, CAP theorem, eventual consistency, and distributed transactions.",
					Resources:   []string{"Designing Data-Intensive Applications (Kleppmann)", "Distributed Systems (van Steen)", "https://jepsen.io"},
				},
				{
					Name:        "Databases",
					Description: "Persistent data storage and retrieval systems. Covers relational (SQL), NoSQL (document, key-value, graph), ACID properties, indexing, and query optimization.",
					Resources:   []string{"Database Internals (Petrov)", "Use The Index, Luke!", "https://db-engines.com"},
				},
			},
		},
		{
			Name:        "Languages",
			Description: "Programming language paradigms and their implementations. Understanding paradigms helps choose the right tool and write idiomatic code.",
			Resources:   []string{"Programming Language Pragmatics (Scott)", "Concepts of Programming Languages (Sebesta)", "https://rosettacode.org"},
			Children: []*Node{
				{
					Name:        "Compiled",
					Description: "Languages translated to machine code before execution (C, C++, Go, Rust). Offer high performance, static type checking, and direct hardware access.",
					Resources:   []string{"The Go Programming Language (Donovan)", "Programming Rust (Blandy)", "https://gobyexample.com"},
				},
				{
					Name:        "Interpreted",
					Description: "Languages executed line-by-line by an interpreter (Python, Ruby, JavaScript). Offer rapid development, dynamic typing, and interactive REPLs.",
					Resources:   []string{"Fluent Python (Ramalho)", "Eloquent JavaScript (Haverbeke)", "https://replit.com"},
				},
				{
					Name:        "Functional",
					Description: "Languages emphasizing pure functions, immutability, and composition (Haskell, Erlang, Clojure, Elixir). Enable reasoning about correctness and concurrency.",
					Resources:   []string{"Learn You a Haskell", "Structure and Interpretation of Computer Programs", "https://exercism.org/tracks/haskell"},
				},
				{
					Name:        "Logic",
					Description: "Languages based on formal logic and declarative rules (Prolog, Datalog, miniKanren). Programs define what to compute, not how. Used in AI, databases, and verification.",
					Resources:   []string{"The Art of Prolog (Sterling)", "Reasoned Schemer (Friedman)", "https://swish.swi-prolog.org"},
				},
			},
		},
	}

	setParents(tree, nil)
	return tree
}

// Styles
var (
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	selectedStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	categoryStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	breadcrumbStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Italic(true)
	detailStyle     = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	resourceStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
)

type model struct {
	roots      []*Node
	cursor     int
	showDetail bool
	jsonMode   bool
}

func initialModel(jsonMode bool) model {
	roots := BuildKnowledgeTree()
	return model{
		roots:    roots,
		jsonMode: jsonMode,
	}
}

func (m model) visible() []*Node {
	return FlattenVisible(m.roots)
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		vis := m.visible()
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(vis)-1 {
				m.cursor++
			}
			m.showDetail = false
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
			m.showDetail = false
		case "enter", " ":
			if m.cursor < len(vis) {
				node := vis[m.cursor]
				if node.IsLeaf() {
					m.showDetail = !m.showDetail
				} else {
					node.Expanded = !node.Expanded
					// If collapsing, ensure cursor doesn't point past visible
					newVis := m.visible()
					if m.cursor >= len(newVis) {
						m.cursor = len(newVis) - 1
					}
				}
			}
		case "l", "right":
			// Expand or show detail
			if m.cursor < len(vis) {
				node := vis[m.cursor]
				if node.IsLeaf() {
					m.showDetail = true
				} else if !node.Expanded {
					node.Expanded = true
				}
			}
		case "h", "left":
			// Collapse or go to parent
			if m.cursor < len(vis) {
				node := vis[m.cursor]
				if !node.IsLeaf() && node.Expanded {
					node.Expanded = false
					newVis := m.visible()
					if m.cursor >= len(newVis) {
						m.cursor = len(newVis) - 1
					}
				} else if node.Parent != nil {
					// Move cursor to parent
					newVis := m.visible()
					for i, n := range newVis {
						if n == node.Parent {
							m.cursor = i
							break
						}
					}
				}
				m.showDetail = false
			}
		case "?":
			m.showDetail = !m.showDetail
		case "esc":
			m.showDetail = false
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.jsonMode {
		return m.jsonView()
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("Knowledge Explorer"))
	b.WriteString("\n")

	vis := m.visible()

	// Breadcrumb
	if m.cursor < len(vis) {
		bc := vis[m.cursor].Breadcrumb()
		b.WriteString(breadcrumbStyle.Render(bc))
	}
	b.WriteString("\n\n")

	// Tree view
	for i, node := range vis {
		depth := node.Depth()
		indent := strings.Repeat("  ", depth)

		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		prefix := ""
		if !node.IsLeaf() {
			if node.Expanded {
				prefix = "▼ "
			} else {
				prefix = "▶ "
			}
		} else {
			prefix = "  "
		}

		name := node.Name
		if i == m.cursor {
			name = selectedStyle.Render(node.Name)
		} else if !node.IsLeaf() {
			name = categoryStyle.Render(node.Name)
		}

		b.WriteString(fmt.Sprintf("%s%s%s%s\n", cursor, indent, prefix, name))
	}

	// Detail view
	if m.showDetail && m.cursor < len(vis) {
		node := vis[m.cursor]
		var detail strings.Builder
		detail.WriteString(titleStyle.Render(node.Name))
		detail.WriteString("\n\n")
		detail.WriteString(node.Description)
		if len(node.Resources) > 0 {
			detail.WriteString("\n\n")
			detail.WriteString(categoryStyle.Render("Resources:"))
			for _, r := range node.Resources {
				detail.WriteString("\n  ")
				detail.WriteString(resourceStyle.Render("• " + r))
			}
		}
		b.WriteString("\n")
		b.WriteString(detailStyle.Render(detail.String()))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("j/k: navigate  enter: expand/detail  h/l: collapse/expand  ?: toggle detail  q: quit"))

	return b.String()
}

type jsonNode struct {
	Name     string     `json:"name"`
	Children []jsonNode `json:"children,omitempty"`
}

func toJSON(nodes []*Node) []jsonNode {
	var result []jsonNode
	for _, n := range nodes {
		jn := jsonNode{Name: n.Name}
		if len(n.Children) > 0 {
			jn.Children = toJSON(n.Children)
		}
		result = append(result, jn)
	}
	return result
}

func (m model) jsonView() string {
	data := toJSON(m.roots)
	out, _ := json.MarshalIndent(data, "", "  ")
	return string(out)
}

func main() {
	jsonMode := false
	for _, arg := range os.Args[1:] {
		if arg == "--json" {
			jsonMode = true
		}
	}

	if jsonMode {
		m := initialModel(true)
		fmt.Println(m.View())
		return
	}

	m := initialModel(false)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
