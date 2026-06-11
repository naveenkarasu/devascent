fn min_stack_ops(operations: Vec<Vec<J>>) -> Vec<J> {
    let mut stack: Vec<i64> = Vec::new();
    let mut mins: Vec<i64> = Vec::new();
    let mut out: Vec<J> = Vec::new();
    for op in &operations {
        let name = match &op[0] {
            J::Str(s) => s.as_str(),
            _ => "",
        };
        match name {
            "push" => {
                if let J::Int(v) = &op[1] {
                    let v = *v;
                    stack.push(v);
                    let nm = if mins.is_empty() { v } else { v.min(*mins.last().unwrap()) };
                    mins.push(nm);
                }
                out.push(J::Null);
            }
            "pop" => {
                stack.pop();
                mins.pop();
                out.push(J::Null);
            }
            "top" => {
                out.push(J::Int(*stack.last().unwrap()));
            }
            _ => {
                // getMin
                out.push(J::Int(*mins.last().unwrap()));
            }
        }
    }
    out
}
