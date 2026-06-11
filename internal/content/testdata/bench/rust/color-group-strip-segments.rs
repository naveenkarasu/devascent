use std::collections::HashMap;

fn count_segments(groups: Vec<Vec<i64>>, strip: Vec<i64>) -> i64 {
    let mut color_to_group: HashMap<i64, usize> = HashMap::new();
    for (idx, group) in groups.iter().enumerate() {
        for &color in group {
            color_to_group.insert(color, idx);
        }
    }
    if strip.is_empty() {
        return 0;
    }
    let mut segments = 1i64;
    for i in 1..strip.len() {
        if color_to_group.get(&strip[i - 1]) != color_to_group.get(&strip[i]) {
            segments += 1;
        }
    }
    segments
}
