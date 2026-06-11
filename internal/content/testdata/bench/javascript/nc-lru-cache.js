function lru_cache_ops(capacity, operations) {
    const cache = new Map();
    const out = [];
    for (const op of operations) {
        if (op[0] === "get") {
            const key = op[1];
            if (!cache.has(key)) {
                out.push(-1);
            } else {
                const val = cache.get(key);
                cache.delete(key);
                cache.set(key, val);
                out.push(val);
            }
        } else {
            // put
            const key = op[1], value = op[2];
            if (cache.has(key)) cache.delete(key);
            cache.set(key, value);
            if (cache.size > capacity) {
                const oldest = cache.keys().next().value;
                cache.delete(oldest);
            }
            out.push(null);
        }
    }
    return out;
}
