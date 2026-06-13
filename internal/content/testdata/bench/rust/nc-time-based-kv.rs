use std::collections::HashMap;

fn time_map_ops(operations: Vec<Vec<J>>) -> Vec<J> {
    // key -> Vec<(timestamp, value)> with strictly increasing timestamps (append order).
    let mut store: HashMap<String, Vec<(i64, String)>> = HashMap::new();
    let mut out: Vec<J> = Vec::new();

    for op in &operations {
        let name = match &op[0] {
            J::Str(s) => s.as_str(),
            _ => "",
        };
        match name {
            "set" => {
                let key = if let J::Str(s) = &op[1] { s.clone() } else { String::new() };
                let value = if let J::Str(s) = &op[2] { s.clone() } else { String::new() };
                let ts = if let J::Int(v) = &op[3] { *v } else { 0 };
                store.entry(key).or_insert_with(Vec::new).push((ts, value));
                out.push(J::Null);
            }
            _ => {
                // get
                let key = if let J::Str(s) = &op[1] { s.clone() } else { String::new() };
                let ts = if let J::Int(v) = &op[2] { *v } else { 0 };
                let mut result = String::new();
                if let Some(list) = store.get(&key) {
                    // bisect_right: largest index with timestamp <= ts
                    let mut i = list.len();
                    while i > 0 && list[i - 1].0 > ts {
                        i -= 1;
                    }
                    if i > 0 {
                        result = list[i - 1].1.clone();
                    }
                }
                out.push(J::Str(result));
            }
        }
    }
    out
}
