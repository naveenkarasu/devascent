function two_sum(nums, target) {
  const seen = {};
  for (let i = 0; i < nums.length; i++) {
    const need = target - nums[i];
    if (seen[need] !== undefined) {
      return [seen[need], i];
    }
    seen[nums[i]] = i;
  }
  return [];
}
