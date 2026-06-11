fn lru_cache_ops(capacity: i64, operations: Vec<Vec<J>>) -> Vec<J> {
    // Maintain entries in recency order: index 0 = least-recently-used, last = most-recent.
    let mut entries: Vec<(i64, i64)> = Vec::new();
    let cap = capacity;
    let mut out: Vec<J> = Vec::new();

    for op in &operations {
        let name = match &op[0] {
            J::Str(s) => s.as_str(),
            _ => "",
        };
        match name {
            "get" => {
                let key = if let J::Int(v) = &op[1] { *v } else { 0 };
                if let Some(pos) = entries.iter().position(|&(k, _)| k == key) {
                    let entry = entries.remove(pos);
                    let val = entry.1;
                    entries.push(entry);
                    out.push(J::Int(val));
                } else {
                    out.push(J::Int(-1));
                }
            }
            _ => {
                // put
                let key = if let J::Int(v) = &op[1] { *v } else { 0 };
                let value = if let J::Int(v) = &op[2] { *v } else { 0 };
                if let Some(pos) = entries.iter().position(|&(k, _)| k == key) {
                    entries.remove(pos);
                }
                entries.push((key, value));
                if entries.len() as i64 > cap {
                    entries.remove(0);
                }
                out.push(J::Null);
            }
        }
    }
    out
}
