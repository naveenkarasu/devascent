function insert_interval(intervals: number[][], newInterval: number[]): number[][] {
    const res: number[][] = [];
    let i = 0;
    const n = intervals.length;
    let s = newInterval[0], e = newInterval[1];

    while (i < n && intervals[i][1] < s) {
        res.push(intervals[i]);
        i++;
    }
    while (i < n && intervals[i][0] <= e) {
        s = Math.min(s, intervals[i][0]);
        e = Math.max(e, intervals[i][1]);
        i++;
    }
    res.push([s, e]);
    while (i < n) {
        res.push(intervals[i]);
        i++;
    }
    return res;
}
