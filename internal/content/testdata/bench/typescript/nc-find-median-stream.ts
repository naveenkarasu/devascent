class MedianFinder {
    private sorted: number[];
    constructor() {
        this.sorted = [];
    }
    add_num(v: number): void {
        let lo = 0, hi = this.sorted.length;
        while (lo < hi) {
            const mid = Math.floor((lo + hi) / 2);
            if (this.sorted[mid] < v) lo = mid + 1;
            else hi = mid;
        }
        this.sorted.splice(lo, 0, v);
    }
    find_median(): number {
        const len = this.sorted.length;
        if (len % 2 === 1) return this.sorted[Math.floor(len / 2)];
        return (this.sorted[len / 2 - 1] + this.sorted[len / 2]) / 2;
    }
}

function median_stream_ops(operations: any[][]): any[] {
    const mf = new MedianFinder();
    const out: any[] = [];
    for (const op of operations) {
        if (op[0] === "addNum") {
            mf.add_num(op[1]);
            out.push(null);
        } else {
            out.push(mf.find_median());
        }
    }
    return out;
}
