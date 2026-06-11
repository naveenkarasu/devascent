function next_permutation(nums: number[]): number[] {
    const n = nums.length;
    let pivot = -1;
    for (let i = n - 2; i >= 0; i--) {
        if (nums[i] < nums[i + 1]) {
            pivot = i;
            break;
        }
    }
    if (pivot === -1) {
        nums.reverse();
        return nums;
    }
    for (let r = n - 1; r > pivot; r--) {
        if (nums[r] > nums[pivot]) {
            const tmp = nums[pivot];
            nums[pivot] = nums[r];
            nums[r] = tmp;
            break;
        }
    }
    // Reverse suffix from pivot+1
    let left = pivot + 1, right = n - 1;
    while (left < right) {
        const tmp = nums[left];
        nums[left] = nums[right];
        nums[right] = tmp;
        left++;
        right--;
    }
    return nums;
}
