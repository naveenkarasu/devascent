fn generate_parenthesis(n: i64) -> Vec<String> {
    fn backtrack(current: String, open_count: i64, close_count: i64, n: i64, result: &mut Vec<String>) {
        if current.len() as i64 == 2 * n {
            result.push(current);
            return;
        }
        if open_count < n {
            let mut next = current.clone();
            next.push('(');
            backtrack(next, open_count + 1, close_count, n, result);
        }
        if close_count < open_count {
            let mut next = current.clone();
            next.push(')');
            backtrack(next, open_count, close_count + 1, n, result);
        }
    }
    let mut result: Vec<String> = Vec::new();
    backtrack(String::new(), 0, 0, n, &mut result);
    result.sort();
    result
}
