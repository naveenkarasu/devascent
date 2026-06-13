function min_max_pages(pages: number[], num_students: number): number {
    if (num_students > pages.length) return -1;

    function can_allocate(max_pages: number): boolean {
        let students = 1;
        let current = 0;
        for (const p of pages) {
            if (p > max_pages) return false;
            if (current + p > max_pages) {
                students++;
                current = p;
            } else {
                current += p;
            }
        }
        return students <= num_students;
    }

    let lo = Math.max(...pages);
    let hi = pages.reduce((a, b) => a + b, 0);
    let result = hi;
    while (lo <= hi) {
        const mid = Math.floor((lo + hi) / 2);
        if (can_allocate(mid)) {
            result = mid;
            hi = mid - 1;
        } else {
            lo = mid + 1;
        }
    }
    return result;
}
