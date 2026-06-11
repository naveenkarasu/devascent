fn best_team_coverage(topics: Vec<Vec<i64>>) -> Vec<i64> {
    let n = topics.len();
    let k = topics[0].len();
    let mut max_topics: i64 = 0;
    let mut best_count: i64 = 0;
    for i in 0..n {
        for j in (i + 1)..n {
            let mut covered: i64 = 0;
            for t in 0..k {
                if topics[i][t] == 1 || topics[j][t] == 1 {
                    covered += 1;
                }
            }
            if covered > max_topics {
                max_topics = covered;
                best_count = 1;
            } else if covered == max_topics {
                best_count += 1;
            }
        }
    }
    vec![max_topics, best_count]
}
