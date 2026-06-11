function running_medians(stream) {
    // lo: max-heap (as negated min-heap), hi: min-heap
    // Using sorted insertion for simplicity
    class MinHeap {
        constructor(cmp) {
            this.data = [];
            this.cmp = cmp || ((a, b) => a - b);
        }
        push(val) {
            this.data.push(val);
            this._bubbleUp(this.data.length - 1);
        }
        pop() {
            const top = this.data[0];
            const last = this.data.pop();
            if (this.data.length > 0) {
                this.data[0] = last;
                this._sinkDown(0);
            }
            return top;
        }
        peek() { return this.data[0]; }
        size() { return this.data.length; }
        _bubbleUp(i) {
            while (i > 0) {
                const parent = Math.floor((i - 1) / 2);
                if (this.cmp(this.data[i], this.data[parent]) < 0) {
                    [this.data[i], this.data[parent]] = [this.data[parent], this.data[i]];
                    i = parent;
                } else break;
            }
        }
        _sinkDown(i) {
            const n = this.data.length;
            while (true) {
                let smallest = i;
                const l = 2*i+1, r = 2*i+2;
                if (l < n && this.cmp(this.data[l], this.data[smallest]) < 0) smallest = l;
                if (r < n && this.cmp(this.data[r], this.data[smallest]) < 0) smallest = r;
                if (smallest !== i) {
                    [this.data[i], this.data[smallest]] = [this.data[smallest], this.data[i]];
                    i = smallest;
                } else break;
            }
        }
    }

    const lo = new MinHeap((a, b) => b - a); // max-heap
    const hi = new MinHeap((a, b) => a - b); // min-heap
    const result = [];

    for (const num of stream) {
        lo.push(num);
        hi.push(lo.pop());
        if (hi.size() > lo.size()) {
            lo.push(hi.pop());
        }
        if (lo.size() === hi.size()) {
            result.push((lo.peek() + hi.peek()) / 2.0);
        } else {
            result.push(lo.peek() * 1.0);
        }
    }
    return result;
}
