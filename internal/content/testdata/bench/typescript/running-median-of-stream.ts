function running_medians(stream: number[]): number[] {
    // lo: max-heap (stored as negatives in min-heap via sorted array)
    // hi: min-heap
    // Using sorted arrays as simple heap simulation
    const lo: number[] = []; // max-heap: store negatives, sorted ascending
    const hi: number[] = []; // min-heap: sorted ascending

    function pushMaxHeap(val: number): void {
        // insert -val into lo (min-heap of negatives)
        lo.push(-val);
        lo.sort((a, b) => a - b);
    }
    function popMaxHeap(): number {
        return -lo.shift()!;
    }
    function peekMaxHeap(): number {
        return -lo[0];
    }
    function pushMinHeap(val: number): void {
        hi.push(val);
        hi.sort((a, b) => a - b);
    }
    function popMinHeap(): number {
        return hi.shift()!;
    }
    function peekMinHeap(): number {
        return hi[0];
    }

    const result: number[] = [];
    for (const num of stream) {
        pushMaxHeap(num);
        pushMinHeap(popMaxHeap());
        if (hi.length > lo.length) {
            pushMaxHeap(popMinHeap());
        }
        if (lo.length === hi.length) {
            result.push((peekMaxHeap() + peekMinHeap()) / 2.0);
        } else {
            result.push(peekMaxHeap());
        }
    }
    return result;
}
