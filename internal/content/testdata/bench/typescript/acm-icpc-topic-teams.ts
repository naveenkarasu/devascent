function best_team_coverage(topics: number[][]): number[] {
    const n = topics.length;
    const k = topics[0].length;
    let max_topics = 0;
    let best_count = 0;
    for (let i = 0; i < n; i++) {
        for (let j = i + 1; j < n; j++) {
            let covered = 0;
            for (let t = 0; t < k; t++) {
                if (topics[i][t] === 1 || topics[j][t] === 1) covered++;
            }
            if (covered > max_topics) {
                max_topics = covered;
                best_count = 1;
            } else if (covered === max_topics) {
                best_count++;
            }
        }
    }
    return [max_topics, best_count];
}
