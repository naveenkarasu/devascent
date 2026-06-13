use std::collections::HashMap;

fn match_colored_houses(left_colors: Vec<i64>, right_colors: Vec<i64>) -> Vec<i64> {
    let mut color_to_pos: HashMap<i64, i64> = HashMap::new();
    for (i, c) in right_colors.iter().enumerate() {
        color_to_pos.insert(*c, (i as i64) + 1);
    }
    left_colors.iter().map(|c| color_to_pos[c]).collect()
}
