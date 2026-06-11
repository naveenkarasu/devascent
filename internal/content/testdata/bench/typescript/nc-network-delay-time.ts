function network_delay_time(times: number[][], n: number, k: number): number {
    const graph: {[key: number]: [number, number][]} = {};
    for (const [u, v, w] of times) {
        if (!(u in graph)) graph[u] = [];
        graph[u].push([w, v]);
    }

    const dist: {[key: number]: number} = {};
    const heap: [number, number][] = [[0, k]];

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

    while (heap.length > 0) {
        const [d, node] = heapPop();
        if (node in dist) continue;
        dist[node] = d;
        if (graph[node]) {
            for (const [weight, nb] of graph[node]) {
                if (!(nb in dist)) {
                    heapPush([d + weight, nb]);
                }
            }
        }
    }

    if (Object.keys(dist).length !== n) return -1;
    return Math.max(...Object.values(dist));
}
