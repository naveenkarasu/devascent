use std::collections::VecDeque;

fn walls_and_gates(rooms: Vec<Vec<i64>>) -> Vec<Vec<i64>> {
    let mut rooms = rooms;
    if rooms.is_empty() || rooms[0].is_empty() {
        return rooms;
    }
    const INF: i64 = 2147483647;
    let m = rooms.len();
    let n = rooms[0].len();
    let mut queue: VecDeque<(usize, usize)> = VecDeque::new();
    for i in 0..m {
        for j in 0..n {
            if rooms[i][j] == 0 {
                queue.push_back((i, j));
            }
        }
    }
    let dirs: [(i64, i64); 4] = [(1, 0), (-1, 0), (0, 1), (0, -1)];
    while let Some((i, j)) = queue.pop_front() {
        for &(di, dj) in dirs.iter() {
            let ni = i as i64 + di;
            let nj = j as i64 + dj;
            if ni >= 0 && ni < m as i64 && nj >= 0 && nj < n as i64 {
                let (niu, nju) = (ni as usize, nj as usize);
                if rooms[niu][nju] == INF {
                    rooms[niu][nju] = rooms[i][j] + 1;
                    queue.push_back((niu, nju));
                }
            }
        }
    }
    rooms
}
