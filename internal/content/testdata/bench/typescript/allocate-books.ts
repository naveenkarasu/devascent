function allocate_books(pages: number[], students: number): number {
    const n = pages.length;
    if (students > n) return -1;
    function canAllocate(maxPages: number): boolean {
        let count = 1, total = 0;
        for (const p of pages) {
            if (p > maxPages) return false;
            if (total + p > maxPages) {
                count++;
                total = p;
            } else {
                total += p;
            }
        }
        return count <= students;
    }
    let lo = Math.max(...pages), hi = pages.reduce((a, b) => a + b, 0);
    let ans = hi;
    while (lo <= hi) {
        const mid = Math.floor((lo + hi) / 2);
        if (canAllocate(mid)) {
            ans = mid;
            hi = mid - 1;
        } else {
            lo = mid + 1;
        }
    }
    return ans;
}
