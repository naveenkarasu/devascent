fn int_to_roman(num: i64) -> String {
    let mapping: [(i64, &str); 13] = [
        (1000, "M"), (900, "CM"), (500, "D"), (400, "CD"),
        (100, "C"), (90, "XC"), (50, "L"), (40, "XL"),
        (10, "X"), (9, "IX"), (5, "V"), (4, "IV"), (1, "I"),
    ];
    let mut num = num;
    let mut parts = String::new();
    for &(value, numeral) in mapping.iter() {
        while num >= value {
            parts.push_str(numeral);
            num -= value;
        }
    }
    parts
}
