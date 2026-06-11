fn alien_order(words: Vec<String>) -> String {
    use std::collections::{BTreeSet, BinaryHeap, HashMap};
    use std::cmp::Reverse;

    // adjacency: char -> set of chars that come after
    let mut adj: HashMap<char, BTreeSet<char>> = HashMap::new();
    let mut indeg: HashMap<char, i64> = HashMap::new();
    for w in &words {
        for c in w.chars() {
            indeg.entry(c).or_insert(0);
            adj.entry(c).or_default();
        }
    }

    let wv: Vec<Vec<char>> = words.iter().map(|w| w.chars().collect()).collect();
    for i in 0..wv.len().saturating_sub(1) {
        let a = &wv[i];
        let b = &wv[i + 1];
        let m = a.len().min(b.len());
        if a.len() > b.len() && a[..m] == b[..m] {
            return String::new();
        }
        for j in 0..m {
            if a[j] != b[j] {
                let set = adj.get_mut(&a[j]).unwrap();
                if !set.contains(&b[j]) {
                    set.insert(b[j]);
                    *indeg.get_mut(&b[j]).unwrap() += 1;
                }
                break;
            }
        }
    }

    let mut heap: BinaryHeap<Reverse<char>> = BinaryHeap::new();
    for (&c, &d) in &indeg {
        if d == 0 {
            heap.push(Reverse(c));
        }
    }
    let mut res: Vec<char> = Vec::new();
    while let Some(Reverse(c)) = heap.pop() {
        res.push(c);
        let nexts: Vec<char> = adj.get(&c).map(|s| s.iter().cloned().collect()).unwrap_or_default();
        for nxt in nexts {
            let d = indeg.get_mut(&nxt).unwrap();
            *d -= 1;
            if *d == 0 {
                heap.push(Reverse(nxt));
            }
        }
    }
    if res.len() != indeg.len() {
        return String::new();
    }
    res.into_iter().collect()
}
