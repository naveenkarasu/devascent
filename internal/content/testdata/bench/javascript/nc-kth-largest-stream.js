function kth_largest_stream_ops(k, nums, operations) {
    // Min-heap implementation
    class MinHeap {
        constructor() { this.h = []; }
        push(v) {
            this.h.push(v);
            this._bubbleUp(this.h.length - 1);
        }
        pop() {
            const top = this.h[0];
            const last = this.h.pop();
            if (this.h.length > 0) {
                this.h[0] = last;
                this._sinkDown(0);
            }
            return top;
        }
        peek() { return this.h[0]; }
        size() { return this.h.length; }
        _bubbleUp(i) {
            while (i > 0) {
                const p = Math.floor((i - 1) / 2);
                if (this.h[p] <= this.h[i]) break;
                [this.h[p], this.h[i]] = [this.h[i], this.h[p]];
                i = p;
            }
        }
        _sinkDown(i) {
            const n = this.h.length;
            while (true) {
                let smallest = i;
                const l = 2 * i + 1, r = 2 * i + 2;
                if (l < n && this.h[l] < this.h[smallest]) smallest = l;
                if (r < n && this.h[r] < this.h[smallest]) smallest = r;
                if (smallest === i) break;
                [this.h[smallest], this.h[i]] = [this.h[i], this.h[smallest]];
                i = smallest;
            }
        }
    }

    const heap = new MinHeap();
    for (const n of nums) heap.push(n);
    while (heap.size() > k) heap.pop();

    const out = [];
    for (const op of operations) {
        heap.push(op[1]);
        while (heap.size() > k) heap.pop();
        out.push(heap.peek());
    }
    return out;
}
