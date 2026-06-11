function find_kth_largest(nums: number[], k: number): number {
    // Min heap of size k
    const heap: number[] = [];

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

    function heapPush(arr: number[], val: number): void {
        arr.push(val);
        heapifyUp(arr, arr.length - 1);
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

    for (const n of nums) {
        heapPush(heap, n);
        if (heap.length > k) {
            heapPop(heap);
        }
    }
    return heap[0];
}
