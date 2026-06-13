fn majority_elements_n3(nums: Vec<i64>) -> Vec<i64> {
    let mut cnt1: i64 = 0;
    let mut cnt2: i64 = 0;
    let mut num1: Option<i64> = None;
    let mut num2: Option<i64> = None;
    for &n in &nums {
        if Some(n) == num1 {
            cnt1 += 1;
        } else if Some(n) == num2 {
            cnt2 += 1;
        } else if cnt1 == 0 {
            num1 = Some(n);
            cnt1 = 1;
        } else if cnt2 == 0 {
            num2 = Some(n);
            cnt2 = 1;
        } else {
            cnt1 -= 1;
            cnt2 -= 1;
        }
    }
    let threshold = (nums.len() / 3) as i64;
    let mut result: Vec<i64> = Vec::new();
    if let Some(v1) = num1 {
        let c = nums.iter().filter(|&&x| x == v1).count() as i64;
        if c > threshold {
            result.push(v1);
        }
    }
    if let Some(v2) = num2 {
        if num2 != num1 {
            let c = nums.iter().filter(|&&x| x == v2).count() as i64;
            if c > threshold {
                result.push(v2);
            }
        }
    }
    result.sort();
    result
}
