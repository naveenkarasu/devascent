fn partition(s: String) -> Vec<Vec<String>> {
    let bytes = s.as_bytes();
    let n = bytes.len();
    let mut res: Vec<Vec<String>> = Vec::new();
    let mut current: Vec<String> = Vec::new();

    fn is_palindrome(b: &[u8]) -> bool {
        if b.is_empty() {
            return true;
        }
        let mut i = 0;
        let mut j = b.len() - 1;
        while i < j {
            if b[i] != b[j] {
                return false;
            }
            i += 1;
            j -= 1;
        }
        true
    }

    fn backtrack(
        start: usize,
        bytes: &[u8],
        n: usize,
        current: &mut Vec<String>,
        res: &mut Vec<Vec<String>>,
    ) {
        if start == n {
            res.push(current.clone());
            return;
        }
        for end in (start + 1)..=n {
            let sub = &bytes[start..end];
            if is_palindrome(sub) {
                current.push(String::from_utf8(sub.to_vec()).unwrap());
                backtrack(end, bytes, n, current, res);
                current.pop();
            }
        }
    }

    backtrack(0, bytes, n, &mut current, &mut res);
    res.sort();
    res
}
