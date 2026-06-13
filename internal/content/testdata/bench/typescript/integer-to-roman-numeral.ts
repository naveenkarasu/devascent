function int_to_roman(num: number): string {
    const mapping: [number, string][] = [
        [1000, 'M'], [900, 'CM'], [500, 'D'], [400, 'CD'],
        [100, 'C'], [90, 'XC'], [50, 'L'], [40, 'XL'],
        [10, 'X'], [9, 'IX'], [5, 'V'], [4, 'IV'], [1, 'I']
    ];
    const parts: string[] = [];
    for (const [value, numeral] of mapping) {
        while (num >= value) {
            parts.push(numeral);
            num -= value;
        }
    }
    return parts.join('');
}
