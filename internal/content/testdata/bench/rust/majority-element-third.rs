fn majority_element_ii(nums: Vec<i64>) -> Vec<i64> {
    let mut cnt1 = 0i64;
    let mut cnt2 = 0i64;
    let mut num1 = 0i64;
    let mut num2 = 1i64;
    for &n in &nums {
        if num1 == n {
            cnt1 += 1;
        } else if num2 == n {
            cnt2 += 1;
        } else if cnt1 == 0 {
            num1 = n;
            cnt1 = 1;
        } else if cnt2 == 0 {
            num2 = n;
            cnt2 = 1;
        } else {
            cnt1 -= 1;
            cnt2 -= 1;
        }
    }
    let mut c1 = 0i64;
    let mut c2 = 0i64;
    for &n in &nums {
        if n == num1 {
            c1 += 1;
        } else if n == num2 {
            c2 += 1;
        }
    }
    let mut res: Vec<i64> = Vec::new();
    if c1 > nums.len() as i64 / 3 {
        res.push(num1);
    }
    if c2 > nums.len() as i64 / 3 {
        res.push(num2);
    }
    res.sort();
    res
}
