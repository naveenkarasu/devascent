function best_team_coverage(topics) {
    let n = topics.length;
    let k = topics[0].length;
    let maxTopics = 0;
    let bestCount = 0;
    for (let i = 0; i < n; i++) {
        for (let j = i + 1; j < n; j++) {
            let covered = 0;
            for (let t = 0; t < k; t++) {
                if (topics[i][t] === 1 || topics[j][t] === 1) covered++;
            }
            if (covered > maxTopics) {
                maxTopics = covered;
                bestCount = 1;
            } else if (covered === maxTopics) {
                bestCount++;
            }
        }
    }
    return [maxTopics, bestCount];
}
