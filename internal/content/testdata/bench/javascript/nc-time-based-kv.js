function time_map_ops(operations) {
    const store = new Map();
    const out = [];

    function bisect_right(arr, val) {
        let lo = 0, hi = arr.length;
        while (lo < hi) {
            const mid = Math.floor((lo + hi) / 2);
            if (arr[mid] <= val) lo = mid + 1;
            else hi = mid;
        }
        return lo;
    }

    for (const op of operations) {
        if (op[0] === "set") {
            const key = op[1], value = op[2], timestamp = op[3];
            if (!store.has(key)) store.set(key, { times: [], vals: [] });
            const entry = store.get(key);
            entry.times.push(timestamp);
            entry.vals.push(value);
            out.push(null);
        } else {
            // get
            const key = op[1], timestamp = op[2];
            if (!store.has(key)) {
                out.push("");
            } else {
                const { times, vals } = store.get(key);
                const i = bisect_right(times, timestamp);
                out.push(i > 0 ? vals[i - 1] : "");
            }
        }
    }
    return out;
}
