function majority_element_ii(nums) {
    let cnt1 = 0, cnt2 = 0, num1 = 0, num2 = 1;
    for (const n of nums) {
        if (num1 === n) cnt1++;
        else if (num2 === n) cnt2++;
        else if (cnt1 === 0) { num1 = n; cnt1 = 1; }
        else if (cnt2 === 0) { num2 = n; cnt2 = 1; }
        else { cnt1--; cnt2--; }
    }
    cnt1 = 0; cnt2 = 0;
    for (const n of nums) {
        if (n === num1) cnt1++;
        else if (n === num2) cnt2++;
    }
    const res = [];
    if (cnt1 > Math.floor(nums.length / 3)) res.push(num1);
    if (cnt2 > Math.floor(nums.length / 3)) res.push(num2);
    return res.sort((a, b) => a - b);
}
