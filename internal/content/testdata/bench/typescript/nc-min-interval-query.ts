function min_interval(intervals: number[][], queries: number[]): number[] {
    // Sort intervals by start
    intervals.sort((a, b) => a[0] - b[0]);
    // Pair queries with original indices and sort by query value
    const indexed_queries: [number, number][] = queries.map((q, i) => [i, q]);
    indexed_queries.sort((a, b) => a[1] - b[1]);

    const res: number[] = new Array(queries.length).fill(-1);
    // Min-heap by (size, end) — use sorted array, pop front
    const heap: [number, number][] = []; // [size, end]
    let i = 0;

    for (const [orig_idx, q] of indexed_queries) {
        // Push all intervals whose start <= q
        while (i < intervals.length && intervals[i][0] <= q) {
            const [l, r] = intervals[i];
            const size = r - l + 1;
            // Insert in sorted position by (size, end)
            let lo = 0, hi = heap.length;
            while (lo < hi) {
                const mid = (lo + hi) >> 1;
                if (heap[mid][0] < size || (heap[mid][0] === size && heap[mid][1] < r)) lo = mid + 1;
                else hi = mid;
            }
            heap.splice(lo, 0, [size, r]);
            i++;
        }
        // Remove expired intervals (end < q)
        while (heap.length > 0 && heap[0][1] < q) {
            heap.shift();
        }
        if (heap.length > 0) {
            res[orig_idx] = heap[0][0];
        }
    }
    return res;
}
