function k_closest(points: number[][], k: number): number[][] {
    // Max heap by distance, size k
    const heap: [number, number, number][] = []; // [dist2, x, y]

    function heapifyDown(arr: [number, number, number][], i: number): void {
        const n = arr.length;
        while (true) {
            let largest = i;
            const l = 2 * i + 1, r = 2 * i + 2;
            if (l < n && arr[l][0] > arr[largest][0]) largest = l;
            if (r < n && arr[r][0] > arr[largest][0]) largest = r;
            if (largest === i) break;
            [arr[i], arr[largest]] = [arr[largest], arr[i]];
            i = largest;
        }
    }

    function heapifyUp(arr: [number, number, number][], i: number): void {
        while (i > 0) {
            const parent = Math.floor((i - 1) / 2);
            if (arr[parent][0] >= arr[i][0]) break;
            [arr[i], arr[parent]] = [arr[parent], arr[i]];
            i = parent;
        }
    }

    function heapPush(arr: [number, number, number][], val: [number, number, number]): void {
        arr.push(val);
        heapifyUp(arr, arr.length - 1);
    }

    function heapPop(arr: [number, number, number][]): [number, number, number] {
        const top = arr[0];
        const last = arr.pop()!;
        if (arr.length > 0) {
            arr[0] = last;
            heapifyDown(arr, 0);
        }
        return top;
    }

    for (const [x, y] of points) {
        const dist2 = x * x + y * y;
        heapPush(heap, [dist2, x, y]);
        if (heap.length > k) {
            heapPop(heap);
        }
    }

    // Sort heap result by distance
    const sorted = [...heap].sort((a, b) => a[0] - b[0]);
    return sorted.map(([, x, y]) => [x, y]);
}
