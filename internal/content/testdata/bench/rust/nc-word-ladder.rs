use std::collections::HashSet;
use std::collections::VecDeque;

fn ladder_length(begin_word: String, end_word: String, word_list: Vec<String>) -> i64 {
    let word_set: HashSet<String> = word_list.into_iter().collect();
    if !word_set.contains(&end_word) {
        return 0;
    }
    let mut queue: VecDeque<(String, i64)> = VecDeque::new();
    queue.push_back((begin_word.clone(), 1));
    let mut visited: HashSet<String> = HashSet::new();
    visited.insert(begin_word);
    let alphabet = "abcdefghijklmnopqrstuvwxyz";
    while let Some((word, length)) = queue.pop_front() {
        let chars: Vec<char> = word.chars().collect();
        for i in 0..chars.len() {
            for c in alphabet.chars() {
                if c == chars[i] {
                    // candidate would equal word; still fine to check, but skip equal char to mirror behavior
                }
                let mut cand = chars.clone();
                cand[i] = c;
                let candidate: String = cand.into_iter().collect();
                if candidate == end_word {
                    return length + 1;
                }
                if word_set.contains(&candidate) && !visited.contains(&candidate) {
                    visited.insert(candidate.clone());
                    queue.push_back((candidate, length + 1));
                }
            }
        }
    }
    0
}
