function merge_intervals(intervals: number[][]): number[][] {
    const res: number[][] = [];
    const sorted = intervals.slice().sort((a, b) => a[0] - b[0] || a[1] - b[1]);
    for (const [s, e] of sorted) {
        if (res.length > 0 && s <= res[res.length - 1][1]) {
            res[res.length - 1][1] = Math.max(res[res.length - 1][1], e);
        } else {
            res.push([s, e]);
        }
    }
    return res;
}
