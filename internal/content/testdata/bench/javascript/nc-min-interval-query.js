function min_interval(intervals, queries) {
    intervals.sort((a, b) => a[0] - b[0]);
    const indexed_queries = queries.map((q, i) => [i, q]).sort((a, b) => a[1] - b[1]);
    const res = new Array(queries.length).fill(-1);
    // Min-heap by size, then end
    const heap = []; // [size, end]

    function heapPush(item) {
        heap.push(item);
        let i = heap.length - 1;
        while (i > 0) {
            const parent = (i - 1) >> 1;
            if (heap[parent][0] > heap[i][0] || (heap[parent][0] === heap[i][0] && heap[parent][1] > heap[i][1])) {
                [heap[parent], heap[i]] = [heap[i], heap[parent]];
                i = parent;
            } else break;
        }
    }

    function heapPop() {
        const top = heap[0];
        const last = heap.pop();
        if (heap.length > 0) {
            heap[0] = last;
            let i = 0;
            while (true) {
                const l = 2 * i + 1, r = 2 * i + 2;
                let smallest = i;
                if (l < heap.length && (heap[l][0] < heap[smallest][0] || (heap[l][0] === heap[smallest][0] && heap[l][1] < heap[smallest][1]))) smallest = l;
                if (r < heap.length && (heap[r][0] < heap[smallest][0] || (heap[r][0] === heap[smallest][0] && heap[r][1] < heap[smallest][1]))) smallest = r;
                if (smallest !== i) { [heap[i], heap[smallest]] = [heap[smallest], heap[i]]; i = smallest; }
                else break;
            }
        }
        return top;
    }

    let i = 0;
    for (const [orig_idx, q] of indexed_queries) {
        while (i < intervals.length && intervals[i][0] <= q) {
            const [l, r] = intervals[i];
            heapPush([r - l + 1, r]);
            i++;
        }
        while (heap.length > 0 && heap[0][1] < q) {
            heapPop();
        }
        if (heap.length > 0) {
            res[orig_idx] = heap[0][0];
        }
    }
    return res;
}
