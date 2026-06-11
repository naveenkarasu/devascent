fn zigzag_convert(s: String, num_rows: i64) -> String {
    if num_rows == 1 {
        return s;
    }
    let nr = num_rows as usize;
    let mut rows: Vec<String> = vec![String::new(); nr];
    let mut row: i64 = 0;
    let mut direction: i64 = -1;
    for ch in s.chars() {
        rows[row as usize].push(ch);
        if row == 0 || row == num_rows - 1 {
            direction = -direction;
        }
        row += direction;
    }
    rows.concat()
}
