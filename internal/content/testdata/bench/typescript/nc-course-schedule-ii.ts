function find_order(num_courses: number, prerequisites: number[][]): number[] {
    const in_degree: number[] = new Array(num_courses).fill(0);
    const adj: number[][] = Array.from({length: num_courses}, () => []);
    for (const [a, b] of prerequisites) {
        adj[b].push(a);
        in_degree[a]++;
    }
    // Min-heap via sorted array: insert in order, shift() for pop
    const heap: number[] = [];
    for (let i = 0; i < num_courses; i++) {
        if (in_degree[i] === 0) {
            // Insert sorted
            let pos = heap.length;
            while (pos > 0 && heap[pos - 1] > i) pos--;
            heap.splice(pos, 0, i);
        }
    }
    const order: number[] = [];
    while (heap.length > 0) {
        const node = heap.shift()!;
        order.push(node);
        for (const nei of adj[node]) {
            in_degree[nei]--;
            if (in_degree[nei] === 0) {
                // Insert in sorted position
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
