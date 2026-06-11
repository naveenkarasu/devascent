function last_stone_weight(stones: number[]): number {
    // Max heap using negation with array-based approach
    const heap: number[] = [...stones].map(s => -s);

    function heapifyDown(arr: number[], i: number): void {
        const n = arr.length;
        while (true) {
            let smallest = i;
            const l = 2 * i + 1, r = 2 * i + 2;
            if (l < n && arr[l] < arr[smallest]) smallest = l;
            if (r < n && arr[r] < arr[smallest]) smallest = r;
            if (smallest === i) break;
            [arr[i], arr[smallest]] = [arr[smallest], arr[i]];
            i = smallest;
        }
    }

    function heapifyUp(arr: number[], i: number): void {
        while (i > 0) {
            const parent = Math.floor((i - 1) / 2);
            if (arr[parent] <= arr[i]) break;
            [arr[i], arr[parent]] = [arr[parent], arr[i]];
            i = parent;
        }
    }

    function heapPop(arr: number[]): number {
        const top = arr[0];
        const last = arr.pop()!;
        if (arr.length > 0) {
            arr[0] = last;
            heapifyDown(arr, 0);
        }
        return top;
    }

    function heapPush(arr: number[], val: number): void {
        arr.push(val);
        heapifyUp(arr, arr.length - 1);
    }

    // Build heap
    for (let i = Math.floor(heap.length / 2) - 1; i >= 0; i--) {
        heapifyDown(heap, i);
    }

    while (heap.length > 1) {
        const y = -heapPop(heap);
        const x = -heapPop(heap);
        if (x !== y) {
            heapPush(heap, -(y - x));
        }
    }
    return heap.length > 0 ? -heap[0] : 0;
}
