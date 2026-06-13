function group_anagrams(strs: string[]): string[][] {
    const groups: Record<string, string[]> = {};
    for (const s of strs) {
        const key = s.split('').sort().join('');
        if (!groups[key]) groups[key] = [];
        groups[key].push(s);
    }
    const res: string[][] = Object.values(groups).map(g => g.sort());
    res.sort((a, b) => {
        for (let i = 0; i < Math.min(a.length, b.length); i++) {
            if (a[i] < b[i]) return -1;
            if (a[i] > b[i]) return 1;
        }
        return a.length - b.length;
    });
    return res;
}
