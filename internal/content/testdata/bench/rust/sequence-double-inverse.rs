fn sequence_equation(p: Vec<i64>) -> Vec<i64> {
    let n = p.len();
    let mut inv = vec![0i64; n + 1];
    for i in 0..n {
        inv[p[i] as usize] = i as i64 + 1;
    }
    (1..=n).map(|x| inv[inv[x] as usize]).collect()
}
