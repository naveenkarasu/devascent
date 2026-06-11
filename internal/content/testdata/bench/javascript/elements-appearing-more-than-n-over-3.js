function majority_elements_n3(nums) {
    let cnt1 = 0, cnt2 = 0, num1 = null, num2 = null;
    for (const n of nums) {
        if (n === num1) {
            cnt1++;
        } else if (n === num2) {
            cnt2++;
        } else if (cnt1 === 0) {
            num1 = n; cnt1 = 1;
        } else if (cnt2 === 0) {
            num2 = n; cnt2 = 1;
        } else {
            cnt1--;
            cnt2--;
        }
    }
    const threshold = Math.floor(nums.length / 3);
    const result = [];
    if (nums.filter(x => x === num1).length > threshold) result.push(num1);
    if (num2 !== num1 && nums.filter(x => x === num2).length > threshold) result.push(num2);
    return result.sort((a, b) => a - b);
}
