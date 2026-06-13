function least_interval(tasks, n) {
    const counts = {};
    for (const t of tasks) counts[t] = (counts[t] || 0) + 1;
    const vals = Object.values(counts);
    const maxCount = Math.max(...vals);
    const numMax = vals.filter(v => v === maxCount).length;
    return Math.max(tasks.length, (maxCount - 1) * (n + 1) + numMax);
}
