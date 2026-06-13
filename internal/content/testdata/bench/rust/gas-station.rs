fn can_complete_circuit(gas: Vec<i64>, cost: Vec<i64>) -> i64 {
    let total_gas: i64 = gas.iter().sum();
    let total_cost: i64 = cost.iter().sum();
    if total_gas < total_cost {
        return -1;
    }
    let mut total: i64 = 0;
    let mut start: i64 = 0;
    for i in 0..gas.len() {
        total += gas[i] - cost[i];
        if total < 0 {
            start = i as i64 + 1;
            total = 0;
        }
    }
    start
}
