function least_interval(tasks: string[], n: number): number {
    const counts = new Map<string, number>();
    for (const t of tasks) {
        counts.set(t, (counts.get(t) || 0) + 1);
    }
    const vals = Array.from(counts.values());
    const maxCount = Math.max(...vals);
    const numMax = vals.filter(v => v === maxCount).length;
    return Math.max(tasks.length, (maxCount - 1) * (n + 1) + numMax);
}
