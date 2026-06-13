fn compare_versions(v1: String, v2: String) -> i64 {
    let parts1: Vec<i64> = v1.split('.').map(|x| x.parse::<i64>().unwrap()).collect();
    let parts2: Vec<i64> = v2.split('.').map(|x| x.parse::<i64>().unwrap()).collect();
    let length = parts1.len().max(parts2.len());
    for i in 0..length {
        let a = if i < parts1.len() { parts1[i] } else { 0 };
        let b = if i < parts2.len() { parts2[i] } else { 0 };
        if a > b {
            return 1;
        } else if a < b {
            return -1;
        }
    }
    0
}
