fn longest_increasing_path(matrix: Vec<Vec<i64>>) -> i64 {
    if matrix.is_empty() || matrix[0].is_empty() {
        return 0;
    }
    let m = matrix.len();
    let n = matrix[0].len();
    let mut memo = vec![vec![0i64; n]; m];

    fn dfs(r: usize, c: usize, matrix: &Vec<Vec<i64>>, memo: &mut Vec<Vec<i64>>, m: usize, n: usize) -> i64 {
        if memo[r][c] != 0 {
            return memo[r][c];
        }
        let mut best = 1i64;
        let dirs: [(i64, i64); 4] = [(-1, 0), (1, 0), (0, -1), (0, 1)];
        for &(dr, dc) in dirs.iter() {
            let nr = r as i64 + dr;
            let nc = c as i64 + dc;
            if nr >= 0 && nr < m as i64 && nc >= 0 && nc < n as i64 {
                let (nru, ncu) = (nr as usize, nc as usize);
                if matrix[nru][ncu] > matrix[r][c] {
                    let cand = 1 + dfs(nru, ncu, matrix, memo, m, n);
                    if cand > best {
                        best = cand;
                    }
                }
            }
        }
        memo[r][c] = best;
        best
    }

    let mut ans = 0i64;
    for r in 0..m {
        for c in 0..n {
            let v = dfs(r, c, &matrix, &mut memo, m, n);
            if v > ans {
                ans = v;
            }
        }
    }
    ans
}
