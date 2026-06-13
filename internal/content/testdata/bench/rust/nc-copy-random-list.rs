fn copy_list(arr: Vec<Vec<J>>) -> Vec<Vec<J>> {
    let mut out: Vec<Vec<J>> = Vec::new();
    for pair in &arr {
        let val = match &pair[0] {
            J::Int(v) => J::Int(*v),
            other => other.clone(),
        };
        let rnd = match pair.get(1) {
            Some(J::Int(r)) => J::Int(*r),
            Some(J::Null) => J::Null,
            Some(other) => other.clone(),
            None => J::Null,
        };
        out.push(vec![val, rnd]);
    }
    out
}
