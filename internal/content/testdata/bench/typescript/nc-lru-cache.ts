class LRUCache {
    private cap: number;
    private cache: Map<number, number>;
    constructor(capacity: number) {
        this.cap = capacity;
        this.cache = new Map();
    }
    get(key: number): number {
        if (!this.cache.has(key)) return -1;
        const val = this.cache.get(key)!;
        this.cache.delete(key);
        this.cache.set(key, val);
        return val;
    }
    put(key: number, value: number): void {
        if (this.cache.has(key)) this.cache.delete(key);
        this.cache.set(key, value);
        if (this.cache.size > this.cap) {
            this.cache.delete(this.cache.keys().next().value);
        }
    }
}

function lru_cache_ops(capacity: number, operations: any[][]): any[] {
    const cache = new LRUCache(capacity);
    const out: any[] = [];
    for (const op of operations) {
        if (op[0] === "get") {
            out.push(cache.get(op[1]));
        } else {
            cache.put(op[1], op[2]);
            out.push(null);
        }
    }
    return out;
}
