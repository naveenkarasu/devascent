fn net_salary(gross: i64) -> i64 {
    let g = gross as f64;
    let tax: f64 = if g <= 300000.0 {
        0.0
    } else if g <= 500000.0 {
        (g - 300000.0) * 0.05
    } else if g <= 1000000.0 {
        200000.0 * 0.05 + (g - 500000.0) * 0.20
    } else {
        200000.0 * 0.05 + 500000.0 * 0.20 + (g - 1000000.0) * 0.30
    };
    (g - tax).round() as i64
}
