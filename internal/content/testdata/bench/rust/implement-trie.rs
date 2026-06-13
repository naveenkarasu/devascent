use std::collections::HashMap;

struct Trie {
    children: HashMap<char, Trie>,
    end: bool,
}

impl Trie {
    fn new() -> Self {
        Trie { children: HashMap::new(), end: false }
    }
    fn insert(&mut self, word: &str) {
        let mut node = self;
        for c in word.chars() {
            node = node.children.entry(c).or_insert_with(Trie::new);
        }
        node.end = true;
    }
    fn search(&self, word: &str) -> bool {
        let mut node = self;
        for c in word.chars() {
            match node.children.get(&c) {
                Some(n) => node = n,
                None => return false,
            }
        }
        node.end
    }
    fn starts_with(&self, prefix: &str) -> bool {
        let mut node = self;
        for c in prefix.chars() {
            match node.children.get(&c) {
                Some(n) => node = n,
                None => return false,
            }
        }
        true
    }
}

fn trie_ops(operations: Vec<Vec<String>>) -> Vec<J> {
    let mut t = Trie::new();
    let mut out: Vec<J> = Vec::new();
    for op in &operations {
        let cmd = op[0].as_str();
        let arg = op[1].as_str();
        match cmd {
            "insert" => {
                t.insert(arg);
                out.push(J::Null);
            }
            "search" => {
                out.push(J::Bool(t.search(arg)));
            }
            _ => {
                out.push(J::Bool(t.starts_with(arg)));
            }
        }
    }
    out
}
