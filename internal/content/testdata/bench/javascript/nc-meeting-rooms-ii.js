function min_meeting_rooms(intervals) {
    if (!intervals || intervals.length === 0) return 0;
    const starts = intervals.map(i => i[0]).sort((a, b) => a - b);
    const ends = intervals.map(i => i[1]).sort((a, b) => a - b);
    let rooms = 0, best = 0, s = 0, e = 0;
    while (s < starts.length) {
        if (starts[s] < ends[e]) {
            rooms++;
            s++;
            best = Math.max(best, rooms);
        } else {
            rooms--;
            e++;
        }
    }
    return best;
}
