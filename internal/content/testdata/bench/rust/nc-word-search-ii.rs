fn find_words(board: Vec<Vec<String>>, words: Vec<String>) -> Vec<String> {
    use std::collections::{BTreeSet, HashMap};

    // Trie node
    #[derive(Default)]
    struct Node {
        children: HashMap<char, Node>,
        word: Option<String>,
    }

    let mut root = Node::default();
    for w in &words {
        let mut node = &mut root;
        for c in w.chars() {
            node = node.children.entry(c).or_default();
        }
        node.word = Some(w.clone());
    }

    if board.is_empty() || board[0].is_empty() {
        return Vec::new();
    }
    let rows = board.len();
    let cols = board[0].len();

    // Convert board to Vec<Vec<char>> (each cell is a single-char string)
    let mut grid: Vec<Vec<char>> = board
        .iter()
        .map(|row| row.iter().map(|s| s.chars().next().unwrap()).collect())
        .collect();

    let mut found: BTreeSet<String> = BTreeSet::new();

    fn dfs(
        r: usize,
        c: usize,
        node: &Node,
        grid: &mut Vec<Vec<char>>,
        rows: usize,
        cols: usize,
        found: &mut BTreeSet<String>,
    ) {
        let ch = grid[r][c];
        let nxt = match node.children.get(&ch) {
            Some(n) => n,
            None => return,
        };
        if let Some(w) = &nxt.word {
            found.insert(w.clone());
        }
        grid[r][c] = '#';
        let dirs: [(i64, i64); 4] = [(1, 0), (-1, 0), (0, 1), (0, -1)];
        for (dr, dc) in dirs.iter() {
            let nr = r as i64 + dr;
            let nc = c as i64 + dc;
            if nr >= 0 && nr < rows as i64 && nc >= 0 && nc < cols as i64 {
                let (nru, ncu) = (nr as usize, nc as usize);
                if grid[nru][ncu] != '#' {
                    dfs(nru, ncu, nxt, grid, rows, cols, found);
                }
            }
        }
        grid[r][c] = ch;
    }

    for r in 0..rows {
        for c in 0..cols {
            dfs(r, c, &root, &mut grid, rows, cols, &mut found);
        }
    }
    found.into_iter().collect()
}
