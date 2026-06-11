fn split_odd_even(s: String) -> Vec<String> {
    let chars: Vec<char> = s.chars().collect();
    let even: String = chars.iter().step_by(2).collect();
    let odd: String = chars.iter().skip(1).step_by(2).collect();
    vec![even, odd]
}
