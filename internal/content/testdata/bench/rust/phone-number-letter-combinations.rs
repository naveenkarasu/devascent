fn letter_combinations(digits: String) -> Vec<String> {
    if digits.is_empty() {
        return Vec::new();
    }
    fn map_digit(d: char) -> &'static str {
        match d {
            '2' => "abc",
            '3' => "def",
            '4' => "ghi",
            '5' => "jkl",
            '6' => "mno",
            '7' => "pqrs",
            '8' => "tuv",
            '9' => "wxyz",
            _ => "",
        }
    }
    let mut results: Vec<String> = vec![String::new()];
    for d in digits.chars() {
        let letters = map_digit(d);
        let mut next: Vec<String> = Vec::new();
        for prev in &results {
            for ch in letters.chars() {
                let mut s = prev.clone();
                s.push(ch);
                next.push(s);
            }
        }
        results = next;
    }
    results
}
