fn min_cost_climbing_stairs(cost: Vec<i64>) -> i64 {
    let n = cost.len();
    let mut cost = cost;
    for i in 2..n {
        cost[i] += cost[i - 1].min(cost[i - 2]);
    }
    cost[n - 1].min(cost[n - 2])
}
