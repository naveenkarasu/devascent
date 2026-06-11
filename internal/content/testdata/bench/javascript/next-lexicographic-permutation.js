function next_permutation(nums) {
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
      [nums[pivot], nums[r]] = [nums[r], nums[pivot]];
      break;
    }
  }
  const suffix = nums.splice(pivot + 1).reverse();
  nums.push(...suffix);
  return nums;
}
