function find_order(num_courses, prerequisites) {
    const in_degree = new Array(num_courses).fill(0);
    const adj = Array.from({length: num_courses}, () => []);
    for (const [a, b] of prerequisites) {
        adj[b].push(a);
        in_degree[a]++;
    }
    // Min-heap using sorted insertion (small N, fine for correctness)
    const heap = [];
    for (let i = 0; i < num_courses; i++) {
        if (in_degree[i] === 0) heap.push(i);
    }
    heap.sort((a, b) => a - b);
    const order = [];
    while (heap.length > 0) {
        const node = heap.shift();
        order.push(node);
        for (const nei of adj[node]) {
            in_degree[nei]--;
            if (in_degree[nei] === 0) {
                // Insert in sorted order
                let lo = 0, hi = heap.length;
                while (lo < hi) {
                    const mid = (lo + hi) >> 1;
                    if (heap[mid] < nei) lo = mid + 1;
                    else hi = mid;
                }
                heap.splice(lo, 0, nei);
            }
        }
    }
    return order.length === num_courses ? order : [];
}
