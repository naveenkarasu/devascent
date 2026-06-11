fn min_max_pages(pages: Vec<i64>, num_students: i64) -> i64 {
    if num_students > pages.len() as i64 {
        return -1;
    }
    fn can_allocate(pages: &[i64], max_pages: i64, num_students: i64) -> bool {
        let mut students: i64 = 1;
        let mut current: i64 = 0;
        for &p in pages {
            if p > max_pages {
                return false;
            }
            if current + p > max_pages {
                students += 1;
                current = p;
            } else {
                current += p;
            }
        }
        students <= num_students
    }
    let mut lo = *pages.iter().max().unwrap();
    let mut hi: i64 = pages.iter().sum();
    let mut result = hi;
    while lo <= hi {
        let mid = (lo + hi) / 2;
        if can_allocate(&pages, mid, num_students) {
            result = mid;
            hi = mid - 1;
        } else {
            lo = mid + 1;
        }
    }
    result
}
