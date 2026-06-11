use std::collections::HashMap;

fn dfs(airport: String, graph: &mut HashMap<String, Vec<String>>, result: &mut Vec<String>) {
    while let Some(list) = graph.get_mut(&airport) {
        if list.is_empty() {
            break;
        }
        let next = list.pop().unwrap();
        dfs(next, graph, result);
    }
    result.push(airport);
}

fn find_itinerary(tickets: Vec<Vec<String>>) -> Vec<String> {
    let mut graph: HashMap<String, Vec<String>> = HashMap::new();
    // Sort tickets ascending, then push in reverse so adjacency lists are
    // reverse-sorted and pop() yields the lexicographically smallest dest.
    let mut sorted = tickets.clone();
    sorted.sort();
    for t in sorted.iter().rev() {
        graph.entry(t[0].clone()).or_insert_with(Vec::new).push(t[1].clone());
    }
    let mut result: Vec<String> = Vec::new();
    dfs("JFK".to_string(), &mut graph, &mut result);
    result.reverse();
    result
}
