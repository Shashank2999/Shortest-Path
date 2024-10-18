package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PathRequest struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}

type Response struct {
	ShortestPath []Point `json:"shortestPath"`
}

type Item struct {
	point    Point
	distance int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func findPath(start, end Point) []Point {
	if start == end {
		return []Point{start}
	}

	pq := make(PriorityQueue, 0)
	distances := make(map[Point]int)
	previous := make(map[Point]Point)

	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			distances[Point{x, y}] = int(^uint(0) >> 1)
		}
	}

	distances[start] = 0
	heap.Push(&pq, &Item{point: start, distance: 0})

	directions := []Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	for pq.Len() > 0 {
		currentItem := heap.Pop(&pq).(*Item)
		current := currentItem.point

		if current == end {

			var path []Point
			for p := current; p != (Point{}); p = previous[p] {
				path = append([]Point{p}, path...)
			}
			return path
		}

		for _, dir := range directions {
			next := Point{current.X + dir.X, current.Y + dir.Y}

			if next.X >= 0 && next.X < 20 && next.Y >= 0 && next.Y < 20 {
				newDistance := distances[current] + 1
				if newDistance < distances[next] {
					distances[next] = newDistance
					previous[next] = current
					heap.Push(&pq, &Item{point: next, distance: newDistance})
				}
			}
		}
	}

	return nil
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		enableCors(&w, r)
		return
	}

	enableCors(&w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req PathRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := findPath(req.Start, req.End)

	response := Response{ShortestPath: path}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func enableCors(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	http.HandleFunc("/find-path", pathHandler)
	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
