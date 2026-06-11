fn merge_triplets(triplets: Vec<Vec<i64>>, target: Vec<i64>) -> bool {
    let mut result = vec![0i64, 0, 0];
    for t in &triplets {
        if t[0] <= target[0] && t[1] <= target[1] && t[2] <= target[2] {
            for i in 0..3 {
                if t[i] > result[i] {
                    result[i] = t[i];
                }
            }
        }
    }
    result == target
}
