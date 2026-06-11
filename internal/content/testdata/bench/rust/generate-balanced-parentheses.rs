fn generate_parentheses(n: i64) -> Vec<String> {
    fn backtrack(s: String, open_count: i64, close_count: i64, n: i64, results: &mut Vec<String>) {
        if s.len() as i64 == 2 * n {
            results.push(s);
            return;
        }
        if open_count < n {
            let mut t = s.clone();
            t.push('(');
            backtrack(t, open_count + 1, close_count, n, results);
        }
        if close_count < open_count {
            let mut t = s.clone();
            t.push(')');
            backtrack(t, open_count, close_count + 1, n, results);
        }
    }
    let mut results: Vec<String> = Vec::new();
    backtrack(String::new(), 0, 0, n, &mut results);
    results
}
