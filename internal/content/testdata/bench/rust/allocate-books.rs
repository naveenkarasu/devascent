fn allocate_books(pages: Vec<i64>, students: i64) -> i64 {
    let n = pages.len() as i64;
    if students > n {
        return -1;
    }
    fn can_allocate(pages: &Vec<i64>, students: i64, max_pages: i64) -> bool {
        let mut count: i64 = 1;
        let mut total: i64 = 0;
        for &p in pages {
            if p > max_pages {
                return false;
            }
            if total + p > max_pages {
                count += 1;
                total = p;
            } else {
                total += p;
            }
        }
        count <= students
    }
    let mut lo = *pages.iter().max().unwrap();
    let mut hi: i64 = pages.iter().sum();
    let mut ans = hi;
    while lo <= hi {
        let mid = (lo + hi) / 2;
        if can_allocate(&pages, students, mid) {
            ans = mid;
            hi = mid - 1;
        } else {
            lo = mid + 1;
        }
    }
    ans
}
