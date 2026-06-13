function is_n_straight_hand(hand, group_size) {
    if (hand.length % group_size !== 0) return false;
    const count = new Map();
    for (const card of hand) {
        count.set(card, (count.get(card) || 0) + 1);
    }
    const keys = [...count.keys()].sort((a, b) => a - b);
    for (const card of keys) {
        const need = count.get(card);
        if (need > 0) {
            for (let i = 0; i < group_size; i++) {
                const cur = count.get(card + i) || 0;
                if (cur < need) return false;
                count.set(card + i, cur - need);
            }
        }
    }
    return true;
}
