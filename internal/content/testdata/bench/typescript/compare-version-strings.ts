function compare_versions(v1: string, v2: string): number {
    const parts1 = v1.split('.').map(Number);
    const parts2 = v2.split('.').map(Number);
    const length = Math.max(parts1.length, parts2.length);
    for (let i = 0; i < length; i++) {
        const a = i < parts1.length ? parts1[i] : 0;
        const b = i < parts2.length ? parts2[i] : 0;
        if (a > b) return 1;
        if (a < b) return -1;
    }
    return 0;
}
