fn is_strong_password(s: String) -> bool {
    if s.chars().count() < 8 {
        return false;
    }
    let has_upper = s.chars().any(|c| c.is_uppercase());
    let has_lower = s.chars().any(|c| c.is_lowercase());
    let has_digit = s.chars().any(|c| c.is_ascii_digit());
    let has_special = s.chars().any(|c| !c.is_alphanumeric());
    has_upper && has_lower && has_digit && has_special
}
