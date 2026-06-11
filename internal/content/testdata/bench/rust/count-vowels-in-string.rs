fn count_vowels(s: String) -> i64 {
    let vowels = "aeiouAEIOU";
    let mut count: i64 = 0;
    for ch in s.chars() {
        if vowels.contains(ch) {
            count += 1;
        }
    }
    count
}
