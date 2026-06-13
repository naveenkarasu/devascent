class TimeMap {
    private store: Map<string, [number[], string[]]>;
    constructor() {
        this.store = new Map();
    }
    set(key: string, value: string, timestamp: number): void {
        if (!this.store.has(key)) this.store.set(key, [[], []]);
        const [times, vals] = this.store.get(key)!;
        times.push(timestamp);
        vals.push(value);
    }
    get(key: string, timestamp: number): string {
        if (!this.store.has(key)) return "";
        const [times, vals] = this.store.get(key)!;
        // bisect_right: find first index > timestamp
        let lo = 0, hi = times.length;
        while (lo < hi) {
            const mid = Math.floor((lo + hi) / 2);
            if (times[mid] <= timestamp) lo = mid + 1;
            else hi = mid;
        }
        return lo > 0 ? vals[lo - 1] : "";
    }
}

function time_map_ops(operations: any[][]): any[] {
    const tm = new TimeMap();
    const out: any[] = [];
    for (const op of operations) {
        if (op[0] === "set") {
            tm.set(op[1], op[2], op[3]);
            out.push(null);
        } else {
            out.push(tm.get(op[1], op[2]));
        }
    }
    return out;
}
