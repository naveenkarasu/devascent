function min_cost_connect_points(points: number[][]): number {
    const n = points.length;
    const visited: Set<number> = new Set();
    // Min-heap: [cost, index]
    const heap: [number, number][] = [[0, 0]];
    let total = 0;

    // Simple min-heap operations
    function heapPush(item: [number, number]): void {
        heap.push(item);
        let i = heap.length - 1;
        while (i > 0) {
            const parent = Math.floor((i - 1) / 2);
            if (heap[parent][0] > heap[i][0]) {
                [heap[parent], heap[i]] = [heap[i], heap[parent]];
                i = parent;
            } else break;
        }
    }

    function heapPop(): [number, number] {
        const top = heap[0];
        const last = heap.pop()!;
        if (heap.length > 0) {
            heap[0] = last;
            let i = 0;
            while (true) {
                let smallest = i;
                const l = 2 * i + 1, r = 2 * i + 2;
                if (l < heap.length && heap[l][0] < heap[smallest][0]) smallest = l;
                if (r < heap.length && heap[r][0] < heap[smallest][0]) smallest = r;
                if (smallest !== i) {
                    [heap[smallest], heap[i]] = [heap[i], heap[smallest]];
                    i = smallest;
                } else break;
            }
        }
        return top;
    }

    while (visited.size < n) {
        const [cost, i] = heapPop();
        if (visited.has(i)) continue;
        visited.add(i);
        total += cost;
        const [xi, yi] = points[i];
        for (let j = 0; j < n; j++) {
            if (!visited.has(j)) {
                const [xj, yj] = points[j];
                heapPush([Math.abs(xi - xj) + Math.abs(yi - yj), j]);
            }
        }
    }
    return total;
}
