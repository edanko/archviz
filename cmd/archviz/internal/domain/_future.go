package domain

// // Calculate architectural metrics
// type ArchitectureMetrics struct {
//     CouplingFactor    float64  // Ratio of external to internal links
//     CriticalityScores map[string]int
//     DepthDistribution map[int]int
// }

// func (cg *ComponentGraph) Analyze() ArchitectureMetrics {
//     metrics := ArchitectureMetrics{
//         CriticalityScores: make(map[string]int),
//         DepthDistribution: make(map[int]int),
//     }

//     // Calculate link criticality
//     for _, comp := range cg.idMap {
//         metrics.DepthDistribution[comp.Depth]++
//         for _, link := range comp.Links {
//             metrics.CriticalityScores[link.ID]++
//         }
//     }

//     return metrics
// }

// // Detect circular dependencies
// func (cg *ComponentGraph) FindCycles() [][]string {
//     var cycles [][]string
//     visited := make(map[string]bool)

//     for id := range cg.idMap {
//         if !visited[id] {
//             path := cg.detectCycle(id, make(map[string]bool), visited)
//             if len(path) > 0 {
//                 cycles = append(cycles, path)
//             }
//         }
//     }
//     return cycles
// }

// // Add layout properties to components
// type Layout struct {
//     X, Y       float64
//     Collapsed  bool
//     Visibility int // 0=hidden, 1=minimized, 2=full
// }

// type Component struct {
//     // Existing fields
//     Layout Layout
// }

// // Force-directed layout simulation
// func (cg *ComponentGraph) CalculateLayout(iterations int) {
//     // Implement Barnes-Hut optimization for large graphs
//     for i := 0; i < iterations; i++ {
//         cg.applyRepulsion()
//         cg.applyHierarchyConstraints()
//         cg.applyLinkAttraction()
//     }
// }

// type ArchitectureSnapshot struct {
//     Timestamp  time.Time
//     Components map[string]ComponentState
//     Links      map[string]LinkState
// }

// func (cg *ComponentGraph) TakeSnapshot() ArchitectureSnapshot {
//     snapshot := ArchitectureSnapshot{
//         Timestamp:  time.Now(),
//         Components: make(map[string]ComponentState),
//         Links:      make(map[string]LinkState),
//     }

//     for id, comp := range cg.idMap {
//         snapshot.Components[id] = ComponentState{
//             Type:  comp.Type,
//             Tags:  comp.Tags,
//             Depth: comp.Depth,
//         }
//     }

//     return snapshot
// }

// // Diff two architecture versions
// func DiffSnapshots(a, b ArchitectureSnapshot) ChangeSet {
//     // Implement structural diff algorithm
//     return ChangeSet{}
// }

// type GraphObserver interface {
//     OnComponentSelect(comp *Component)
//     OnLinkHover(link Link)
// }

// func (cg *ComponentGraph) RegisterObserver(obs GraphObserver) {
//     // Implement observer pattern for UI interactions
// }

// // Collapse/expand hierarchy
// func (cg *ComponentGraph) ToggleComponent(id string, collapse bool) {
//     comp := cg.idMap[id]
//     comp.Layout.Collapsed = collapse

//     // Propagate visibility changes to children
//     if collapse {
//         cg.applyToChildren(comp, func(c *Component) {
//             c.Layout.Visibility = 1
//         })
//     }
// }

// type ViewPolicy struct {
//     AllowedTypes    map[ComponentType]bool
//     MaxDepth        int
//     SensitiveTags   map[string]bool
// }

// func (cg *ComponentGraph) ApplyViewPolicy(policy ViewPolicy) *ComponentGraph {
//     filtered := NewComponentGraph()

//     for id, comp := range cg.idMap {
//         if policy.IsAllowed(comp) {
//             filtered.AddComponent(id, comp.Title, comp.Links)
//         }
//     }

//     return filtered
// }

// // Example policy implementation
// func (p ViewPolicy) IsAllowed(comp *Component) bool {
//     if !p.AllowedTypes[comp.Type] {
//         return false
//     }

//     for tag := range comp.Tags {
//         if p.SensitiveTags[tag] {
//             return false
//         }
//     }

//     return comp.Depth <= p.MaxDepth
// }

// func (cg *ComponentGraph) GenerateDocumentation(format string) []byte {
//     switch format {
//     case "markdown":
//         return cg.renderMarkdown()
//     case "plantuml":
//         return cg.renderPlantUML()
//     case "mermaid":
//         return cg.renderMermaid()
//     default:
//         return cg.renderJSON()
//     }
// }

// func (cg *ComponentGraph) renderMermaid() []byte {
//     var sb strings.Builder
//     sb.WriteString("graph TD\n")

//     for _, comp := range cg.idMap {
//         for _, link := range comp.Links {
//             sb.WriteString(fmt.Sprintf(
//                 "    %s[%s] --> %s[%s]\n",
//                 comp.ID, comp.Title,
//                 link.ID, cg.idMap[link.ID].Title,
//             ))
//         }
//     }

//     return []byte(sb.String())
// }

// type ArchitectureRule func(*ComponentGraph) []Violation

// var DefaultRules = []ArchitectureRule{
//     NoDirectDBCallsFromUI,
//     RequiredTagsForProduction,
//     LayeredArchitectureConstraint,
// }

// func NoDirectDBCallsFromUI(g *ComponentGraph) []Violation {
//     var violations []Violation

//     for _, comp := range g.Filter(Filter{Types: map[ComponentType]bool{TypeUI: true}}) {
//         for _, link := range comp.Links {
//             if g.idMap[link.ID].Type == TypeDatabase {
//                 violations = append(violations, Violation{
//                     Component: comp.ID,
//                     Message:   "UI directly accessing database",
//                 })
//             }
//         }
//     }

//     return violations
// }

// // Add caching for frequent queries
// type ComponentGraph struct {
//     // Existing fields
//     queryCache *lru.Cache
// }

// // Memoize complex queries
// func (cg *ComponentGraph) memoizedQuery(key string, fn func() interface{}) interface{} {
//     if val, ok := cg.queryCache.Get(key); ok {
//         return val
//     }

//     result := fn()
//     cg.queryCache.Add(key, result)
//     return result
// }

// // Implement batched updates
// func (cg *ComponentGraph) BatchUpdate(updateFunc func(*TxGraph)) {
//     tx := &TxGraph{
//         pendingAdds:    make([]Component, 0, 100),
//         pendingUpdates: make(map[string]Component),
//         pendingDeletes: make(map[string]bool),
//     }

//     updateFunc(tx)
//     cg.applyTransaction(tx)
// }

// // Anomaly detection for architectural smells
// func (cg *ComponentGraph) DetectAnomalies(modelPath string) []Anomaly {
//     // Convert graph to tensor format
//     features := cg.extractGraphFeatures()

//     // Load trained ML model
//     model := tf.LoadModel(modelPath)

//     // Predict anomalies
//     predictions := model.Predict(features)

//     return cg.interpretPredictions(predictions)
// }

// // Generate architecture recommendations
// func (cg *ComponentGraph) GenerateRecommendations() []Recommendation {
//     // Implement graph pattern matching
//     // and known architecture antipattern detection
//     return []Recommendation{}
// }

// type CollaborationSession struct {
//     cursorPositions map[string]CursorPosition
//     changeBuffer    []patch.Operation
// }

// func (cg *ComponentGraph) HandleWebSocket(conn *websocket.Conn) {
//     // Implement operational transformation (OT)
//     // for concurrent graph editing
//     for {
//         msg := readMessage(conn)
//         cg.applyOTPatch(msg.Patch)
//         broadcastToCollaborators(msg)
//     }
// }
