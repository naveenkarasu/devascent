fn group_anagrams(strs: Vec<String>) -> Vec<Vec<String>> {
    use std::collections::HashMap;
    let mut groups: HashMap<String, Vec<String>> = HashMap::new();
    for s in &strs {
        let mut chars: Vec<char> = s.chars().collect();
        chars.sort();
        let key: String = chars.into_iter().collect();
        groups.entry(key).or_default().push(s.clone());
    }
    let mut res: Vec<Vec<String>> = groups
        .into_values()
        .map(|mut g| {
            g.sort();
            g
        })
        .collect();
    res.sort();
    res
}
