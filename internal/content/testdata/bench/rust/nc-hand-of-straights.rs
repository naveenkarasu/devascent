use std::collections::HashMap;

fn is_n_straight_hand(hand: Vec<i64>, group_size: i64) -> bool {
    if group_size == 0 {
        return true;
    }
    if hand.len() as i64 % group_size != 0 {
        return false;
    }
    let mut count: HashMap<i64, i64> = HashMap::new();
    for &card in &hand {
        *count.entry(card).or_insert(0) += 1;
    }
    let mut keys: Vec<i64> = count.keys().cloned().collect();
    keys.sort();
    for card in keys {
        let cur = *count.get(&card).unwrap_or(&0);
        if cur > 0 {
            let need = cur;
            for i in 0..group_size {
                let have = *count.get(&(card + i)).unwrap_or(&0);
                if have < need {
                    return false;
                }
                *count.entry(card + i).or_insert(0) -= need;
            }
        }
    }
    true
}
