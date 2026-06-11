fn encode_decode(strs: Vec<String>) -> Vec<String> {
    fn encode(strs: &Vec<String>) -> String {
        let mut out = String::new();
        for s in strs {
            // length is in chars (Python len), but we decode by chars too
            out.push_str(&s.chars().count().to_string());
            out.push('#');
            out.push_str(s);
        }
        out
    }
    fn decode(s: &str) -> Vec<String> {
        let chars: Vec<char> = s.chars().collect();
        let mut res: Vec<String> = Vec::new();
        let mut i = 0usize;
        while i < chars.len() {
            let mut j = i;
            while chars[j] != '#' {
                j += 1;
            }
            let len_str: String = chars[i..j].iter().collect();
            let length: usize = len_str.parse().unwrap();
            let word: String = chars[j + 1..j + 1 + length].iter().collect();
            res.push(word);
            i = j + 1 + length;
        }
        res
    }
    let s = encode(&strs);
    decode(&s)
}
