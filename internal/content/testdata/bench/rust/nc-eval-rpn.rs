fn eval_rpn(tokens: Vec<String>) -> i64 {
    let mut stack: Vec<i64> = Vec::new();
    for t in tokens.iter() {
        match t.as_str() {
            "+" => {
                let b = stack.pop().unwrap();
                let a = stack.pop().unwrap();
                stack.push(a + b);
            }
            "-" => {
                let b = stack.pop().unwrap();
                let a = stack.pop().unwrap();
                stack.push(a - b);
            }
            "*" => {
                let b = stack.pop().unwrap();
                let a = stack.pop().unwrap();
                stack.push(a * b);
            }
            "/" => {
                let b = stack.pop().unwrap();
                let a = stack.pop().unwrap();
                stack.push(a / b); // truncates toward zero, matching int(a/b)
            }
            _ => {
                stack.push(t.parse::<i64>().unwrap());
            }
        }
    }
    stack[0]
}
