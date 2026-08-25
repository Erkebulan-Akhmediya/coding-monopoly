# CGC Report

_Generated: 2026-08-25 05:28 UTC_

_Scoped to repository: `/home/student/Desktop/coding-monopoly`_


## God Nodes — Highest Fan-In
_These nodes are called from many places. High fan-in increases risk: a change here affects every caller._

| Kind | Name | File | In-degree |
| --- | --- | --- | --- |
| Function | AdminStartGame | room/room.go | 29 |
| Function | AddOrReconnectPlayer | room/room.go | 22 |
| Function | NewRoom | room/room.go | 17 |
| Function | connectClient | ws/ws_test.go | 16 |
| Function | sendJoin | ws/ws_test.go | 15 |
| Function | setupTestServer | ws/ws_test.go | 13 |
| Function | readUntilType | ws/ws_test.go | 13 |
| Function | sendError | ws/client.go | 13 |
| Function | GetActivePlayerID | room/room.go | 13 |
| Function | writeError | admin/admin.go | 11 |
| Function | NewMessage | ws/message.go | 11 |
| Function | GetRoomID | ws/client.go | 11 |
| Function | SubmitAnswer | room/room.go | 10 |
| Function | GetRoomInstance | ws/hub.go | 10 |
| Function | getBaseHttpUrl | services/serverUrls.ts | 10 |


## Most Complex Functions
_Cyclomatic complexity > 10 is a refactoring candidate._

| Function | File | Cyclomatic Complexity |
| --- | --- | --- |
| handleMessage | services/websocketService.ts | 85 |
| handleMessage | services/adminWebsocketService.ts | 56 |
| classifyEffect | services/soundService.ts | 42 |
| Run | ws/hub.go | 41 |
| main | loadtest/main.go | 39 |
| handleIncomingMessage | ws/client.go | 33 |
| validateProblemInput | services/adminApiService.ts | 32 |
| ChooseLevel | room/room.go | 25 |
| getCellGridStyle | components/BoardView.vue | 24 |
| getCellGridStyle | components/AdminSpectatorView.vue | 24 |
| main | seed/main.go | 21 |
| gradeQuestion | room/room.go | 18 |
| problems | admin/admin.go | 18 |
| TestRoom_JoinOrderAndActivePlayer | room/room_test.go | 18 |
| validate | admin/admin.go | 17 |


## Cross-Module Connections
_Calls that cross package boundaries — review for unexpected coupling._

| Caller | Caller File | Callee | Callee File | Confidence |
| --- | --- | --- | --- | --- |
| handleMessage | services/adminWebsocketService.ts | appendEvent | src/adminStore.ts | INFERRED |
| handleMessage | services/adminWebsocketService.ts | appendEvent | src/adminStore.ts | INFERRED |
| handleMessage | services/adminWebsocketService.ts | appendEvent | src/adminStore.ts | INFERRED |
| handleMessage | services/adminWebsocketService.ts | appendEvent | src/adminStore.ts | INFERRED |
| handleMessage | services/adminWebsocketService.ts | appendEvent | src/adminStore.ts | INFERRED |
| handleMessage | services/adminWebsocketService.ts | appendEvent | src/adminStore.ts | INFERRED |
| handleMessage | services/adminWebsocketService.ts | appendEvent | src/adminStore.ts | INFERRED |
| handleMessage | services/adminWebsocketService.ts | appendEvent | src/adminStore.ts | INFERRED |
| handleMessage | services/adminWebsocketService.ts | appendEvent | src/adminStore.ts | INFERRED |
| handleMessage | services/adminWebsocketService.ts | appendEvent | src/adminStore.ts | INFERRED |
| main | loadtest/main.go | ServeWS | ws/handler.go | INFERRED |
| main | loadtest/main.go | NewHub | ws/hub.go | INFERRED |
| TestDeployMux_SubmitWinsAndStaleTimerDoesNotDoubleResolve | server/deploy_fairness_test.go | GetActivePlayerID | room/room.go | INFERRED |
| TestDeployMux_SubmitWinsAndStaleTimerDoesNotDoubleResolve | server/deploy_fairness_test.go | AdminStartGame | room/room.go | INFERRED |
| TestDeployMux_SubmitWinsAndStaleTimerDoesNotDoubleResolve | server/deploy_fairness_test.go | SetDeadlineDurations | room/room.go | INFERRED |
| TestDeployMux_TimeoutThenLateSubmitResolvesOnce | server/deploy_fairness_test.go | GetActivePlayerID | room/room.go | INFERRED |
| TestDeployMux_TimeoutThenLateSubmitResolvesOnce | server/deploy_fairness_test.go | AdminStartGame | room/room.go | INFERRED |
| TestDeployMux_TimeoutThenLateSubmitResolvesOnce | server/deploy_fairness_test.go | SetDeadlineDurations | room/room.go | INFERRED |
| TestDeployMux_QuestionContentStaysOffSpectatorWire | server/deploy_fairness_test.go | AdminStartGame | room/room.go | INFERRED |
| TestDeployMux_HealthAndEmbeddedAssets | server/deploy_fairness_test.go | NewHub | ws/hub.go | INFERRED |


## Potential Dead Code
_Functions with zero callers (not guaranteed dead — may be entry points or called via reflection)._

| Function | File |
| --- | --- |
| created | src/App.vue |
| data | src/App.vue |
| showBoard | src/App.vue |
| addAnswer | components/AdminQuestionFormModal.vue |
| addOption | components/AdminQuestionFormModal.vue |
| data | components/AdminQuestionFormModal.vue |
| handleCancel | components/AdminQuestionFormModal.vue |
| handleSubmit | components/AdminQuestionFormModal.vue |
| handleTypeChange | components/AdminQuestionFormModal.vue |
| removeAnswer | components/AdminQuestionFormModal.vue |
| removeOption | components/AdminQuestionFormModal.vue |
| cancelDelete | components/AdminQuestionList.vue |
| confirmDelete | components/AdminQuestionList.vue |
| data | components/AdminQuestionList.vue |
| difficultyBadgeClass | components/AdminQuestionList.vue |
| filteredProblems | components/AdminQuestionList.vue |
| formatDate | components/AdminQuestionList.vue |
| mounted | components/AdminQuestionList.vue |
| onFilterChange | components/AdminQuestionList.vue |
| onProblemSaved | components/AdminQuestionList.vue |


## Suggested Cypher Queries
_Copy these into `execute_cypher_query` to explore further._

### Callers of a specific function
```cypher
MATCH (caller)-[:CALLS|HEURISTIC_CALLS]->(fn:Function {name: 'yourFunctionName'})
RETURN caller.name, caller.path LIMIT 20
```

### Class hierarchy for a specific class
```cypher
MATCH path = (c:Class {name: 'YourClass'})-[:INHERITS*]->(parent)
RETURN [n IN nodes(path) | n.name] AS hierarchy
```

### Most-injected Spring beans
```cypher
MATCH ()-[:INJECTS]->(bean:Class)
RETURN bean.name, count(*) AS injection_count
ORDER BY injection_count DESC LIMIT 10
```

### All external library dependencies
```cypher
MATCH (m:MavenModule)-[:USES_LIBRARY]->(lib:ExternalLibrary)
RETURN m.artifact_id, lib.group_id, lib.artifact_id, lib.version
ORDER BY lib.artifact_id
```

### CALLS edges with low confidence (potential mis-resolutions)
```cypher
MATCH (a)-[c:CALLS|HEURISTIC_CALLS]->(b)
WHERE c.confidence_label = 'AMBIGUOUS'
RETURN a.name, b.name, c.resolution_tier, a.path LIMIT 20
```
