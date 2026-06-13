use std::collections::HashMap;

fn detect_squares_ops(operations: Vec<Vec<J>>) -> Vec<J> {
    let mut cnt: HashMap<(i64, i64), i64> = HashMap::new();
    let mut out: Vec<J> = Vec::new();

    // Extract [x, y] from a J::Arr of two ints.
    fn point(j: &J) -> (i64, i64) {
        if let J::Arr(a) = j {
            let x = if let J::Int(v) = &a[0] { *v } else { 0 };
            let y = if let J::Int(v) = &a[1] { *v } else { 0 };
            (x, y)
        } else {
            (0, 0)
        }
    }

    for op in &operations {
        let name = match &op[0] {
            J::Str(s) => s.as_str(),
            _ => "",
        };
        match name {
            "add" => {
                let (x, y) = point(&op[1]);
                *cnt.entry((x, y)).or_insert(0) += 1;
                out.push(J::Null);
            }
            _ => {
                // count
                let (px, py) = point(&op[1]);
                let mut total: i64 = 0;
                let items: Vec<((i64, i64), i64)> =
                    cnt.iter().map(|(&k, &v)| (k, v)).collect();
                for ((x, y), c) in items {
                    if (x - px).abs() == (y - py).abs() && x != px && y != py {
                        let c1 = *cnt.get(&(px, y)).unwrap_or(&0);
                        let c2 = *cnt.get(&(x, py)).unwrap_or(&0);
                        total += c * c1 * c2;
                    }
                }
                out.push(J::Int(total));
            }
        }
    }
    out
}
