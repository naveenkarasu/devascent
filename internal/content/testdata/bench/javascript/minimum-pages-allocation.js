function min_max_pages(pages, num_students) {
    if (num_students > pages.length) return -1;

    function canAllocate(maxPages) {
        let students = 1, current = 0;
        for (const p of pages) {
            if (p > maxPages) return false;
            if (current + p > maxPages) {
                students++;
                current = p;
            } else {
                current += p;
            }
        }
        return students <= num_students;
    }

    let lo = Math.max(...pages), hi = pages.reduce((a, b) => a + b, 0);
    let result = hi;
    while (lo <= hi) {
        const mid = Math.floor((lo + hi) / 2);
        if (canAllocate(mid)) {
            result = mid;
            hi = mid - 1;
        } else {
            lo = mid + 1;
        }
    }
    return result;
}
