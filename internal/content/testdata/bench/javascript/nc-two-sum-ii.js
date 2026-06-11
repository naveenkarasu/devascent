function two_sum_ii(numbers, target) {
    let i = 0, j = numbers.length - 1;
    while (i < j) {
        const total = numbers[i] + numbers[j];
        if (total === target) return [i + 1, j + 1];
        if (total < target) i++;
        else j--;
    }
    return [];
}
