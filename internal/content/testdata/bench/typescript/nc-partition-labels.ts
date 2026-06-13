function partition_labels(s: string): number[] {
    const last: Map<string, number> = new Map();
    for (let i = 0; i < s.length; i++) {
        last.set(s[i], i);
    }
    const result: number[] = [];
    let start = 0, end = 0;
    for (let i = 0; i < s.length; i++) {
        end = Math.max(end, last.get(s[i])!);
        if (i === end) {
            result.push(end - start + 1);
            start = i + 1;
        }
    }
    return result;
}
