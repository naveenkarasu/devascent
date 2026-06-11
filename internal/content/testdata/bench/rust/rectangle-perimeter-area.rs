fn rectangle_perimeter_area(a: i64, b: i64) -> Vec<i64> {
    if a <= 0 || b <= 0 {
        return vec![0, 0];
    }
    vec![(a + b) * 2, a * b]
}
