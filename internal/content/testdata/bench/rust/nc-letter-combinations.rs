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

fn backtrack(index: usize, digits: &Vec<char>, current: String, res: &mut Vec<String>) {
    if index == digits.len() {
        res.push(current);
        return;
    }
    for ch in map_digit(digits[index]).chars() {
        let mut next = current.clone();
        next.push(ch);
        backtrack(index + 1, digits, next, res);
    }
}

fn letter_combinations(digits: String) -> Vec<String> {
    if digits.is_empty() {
        return Vec::new();
    }
    let chars: Vec<char> = digits.chars().collect();
    let mut res: Vec<String> = Vec::new();
    backtrack(0, &chars, String::new(), &mut res);
    res.sort();
    res
}
