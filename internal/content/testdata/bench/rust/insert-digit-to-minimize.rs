fn insert_to_minimize(a: String) -> String {
    let chars: Vec<char> = a.chars().collect();
    if chars[0] == '1' {
        let rest: String = chars[1..].iter().collect();
        return format!("{}0{}", chars[0], rest);
    }
    format!("1{}", a)
}
