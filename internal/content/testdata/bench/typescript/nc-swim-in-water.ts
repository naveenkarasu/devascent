function swim_in_water(grid: number[][]): number {
    const n = grid.length;
    const visited: Set<number> = new Set();
    const heap: [number, number, number][] = [[grid[0][0], 0, 0]];
    visited.add(0);

    function heapPush(item: [number, number, number]): void {
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

    function heapPop(): [number, number, number] {
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
        const [t, i, j] = heapPop();
        if (i === n - 1 && j === n - 1) return t;
        for (const [di, dj] of [[1,0],[-1,0],[0,1],[0,-1]]) {
            const ni = i + di, nj = j + dj;
            if (ni >= 0 && ni < n && nj >= 0 && nj < n && !visited.has(ni * n + nj)) {
                visited.add(ni * n + nj);
                heapPush([Math.max(t, grid[ni][nj]), ni, nj]);
            }
        }
    }
    return -1;
}
