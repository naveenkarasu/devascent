fn roman_to_int(s: String) -> i64 {
    fn double_val(pair: &str) -> Option<i64> {
        match pair {
            "CM" => Some(900),
            "CD" => Some(400),
            "XC" => Some(90),
            "XL" => Some(40),
            "IX" => Some(9),
            "IV" => Some(4),
            _ => None,
        }
    }
    fn single_val(c: char) -> i64 {
        match c {
            'M' => 1000,
            'D' => 500,
            'C' => 100,
            'L' => 50,
            'X' => 10,
            'V' => 5,
            'I' => 1,
            _ => 0,
        }
    }
    let chars: Vec<char> = s.chars().collect();
    let mut total: i64 = 0;
    let mut i = 0usize;
    while i < chars.len() {
        if i < chars.len() - 1 {
            let pair: String = chars[i..i + 2].iter().collect();
            if let Some(v) = double_val(&pair) {
                total += v;
                i += 2;
                continue;
            }
        }
        total += single_val(chars[i]);
        i += 1;
    }
    total
}
