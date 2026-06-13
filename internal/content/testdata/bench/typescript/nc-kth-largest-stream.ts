class KthLargest {
    private k: number;
    private heap: number[];
    constructor(k: number, nums: number[]) {
        this.k = k;
        this.heap = nums.slice().sort((a, b) => a - b);
        while (this.heap.length > k) this.heap.shift();
    }
    add(val: number): number {
        // Insert maintaining sorted order
        let lo = 0, hi = this.heap.length;
        while (lo < hi) {
            const mid = Math.floor((lo + hi) / 2);
            if (this.heap[mid] < val) lo = mid + 1;
            else hi = mid;
        }
        this.heap.splice(lo, 0, val);
        while (this.heap.length > this.k) this.heap.shift();
        return this.heap[0];
    }
}

function kth_largest_stream_ops(k: number, nums: number[], operations: any[][]): any[] {
    const kl = new KthLargest(k, nums);
    const out: any[] = [];
    for (const op of operations) {
        out.push(kl.add(op[1]));
    }
    return out;
}
