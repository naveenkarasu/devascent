function erase_overlap_intervals(intervals) {
    if (!intervals || intervals.length === 0) return 0;
    intervals.sort((a, b) => a[1] - b[1]);
    let end = intervals[0][1];
    let count = 0;
    for (let i = 1; i < intervals.length; i++) {
        if (intervals[i][0] < end) {
            count++;
        } else {
            end = intervals[i][1];
        }
    }
    return count;
}
