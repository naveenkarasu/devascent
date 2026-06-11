fn word_dictionary_ops(operations: Vec<Vec<String>>) -> Vec<J> {
    use std::collections::HashMap;

    #[derive(Default)]
    struct TrieNode {
        children: HashMap<char, TrieNode>,
        is_end: bool,
    }

    fn add_word(root: &mut TrieNode, word: &str) {
        let mut node = root;
        for c in word.chars() {
            node = node.children.entry(c).or_default();
        }
        node.is_end = true;
    }

    fn search(node: &TrieNode, word: &[char], i: usize) -> bool {
        if i == word.len() {
            return node.is_end;
        }
        let c = word[i];
        if c == '.' {
            for child in node.children.values() {
                if search(child, word, i + 1) {
                    return true;
                }
            }
            false
        } else {
            match node.children.get(&c) {
                Some(child) => search(child, word, i + 1),
                None => false,
            }
        }
    }

    let mut root = TrieNode::default();
    let mut out: Vec<J> = Vec::new();
    for op in &operations {
        let name = op[0].as_str();
        let arg = op[1].clone();
        if name == "addWord" {
            add_word(&mut root, &arg);
            out.push(J::Null);
        } else {
            let chars: Vec<char> = arg.chars().collect();
            let res = search(&root, &chars, 0);
            out.push(J::Bool(res));
        }
    }
    out
}
