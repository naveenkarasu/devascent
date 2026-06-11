function median_stream_ops(operations) {
    // MinHeap
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

    // lo is a max-heap implemented as a min-heap with negated values
    const lo = new MinHeap(); // max-heap (negated)
    const hi = new MinHeap(); // min-heap

    const out = [];
    for (const op of operations) {
        if (op[0] === "addNum") {
            const v = op[1];
            lo.push(-v);
            hi.push(-lo.pop());
            if (hi.size() > lo.size()) lo.push(-hi.pop());
            out.push(null);
        } else {
            // findMedian
            if (lo.size() > hi.size()) {
                out.push(-lo.peek());
            } else {
                out.push((-lo.peek() + hi.peek()) / 2);
            }
        }
    }
    return out;
}
